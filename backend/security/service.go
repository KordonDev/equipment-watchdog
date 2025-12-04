package security

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

type userService interface {
	GetUser(string) (*models.User, error)
	AddUser(*models.User) (*models.User, error)
	SaveUser(*models.User) (*models.User, error)
	HasApprovedAndAdminUser() bool
	CheckLogin(username, password string) error
	ChangePassword(ctx context.Context, username, password string) error
	GetUserByCredentialId(string) (*models.User, error)
}

type SessionStore interface {
	getSession(username string) (webauthn.SessionData, error)
	storeSession(username string, sessionData webauthn.SessionData) error
}

type WebAuthNService struct {
	webAuthn     *webauthn.WebAuthn
	jwtService   *JwtService
	userService  userService
	domain       string
	sessionStore SessionStore
}

func NewWebAuthNService(userService userService, origin string, domain string, jwtService *JwtService, db *gorm.DB) (*WebAuthNService, error) {
	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "equipment watchdog", // Display Name for your site
		RPID:          domain,               // Generally the domain name for your site
		RPOrigin:      origin,               // The origin URL for WebAuthn requests
		Timeouts: webauthn.TimeoutsConfig{
			Login: webauthn.TimeoutConfig{
				Enforce:    true,
				Timeout:    2 * time.Minute,
				TimeoutUVD: 2 * time.Minute,
			},
			Registration: webauthn.TimeoutConfig{
				Enforce:    true,
				Timeout:    2 * time.Minute,
				TimeoutUVD: 2 * time.Minute,
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("could not create webAuth: %w", err)
	}

	sessionStore := NewDatebase(db)

	return &WebAuthNService{
		webAuthn:     webAuthn,
		jwtService:   jwtService,
		userService:  userService,
		domain:       domain,
		sessionStore: sessionStore,
	}, nil
}

func (w WebAuthNService) startRegisterNewUser(username string) (*protocol.CredentialCreation, error) {
	if w.userExists(username) {
		return nil, errors.New("user already exists")
	}

	user := &models.User{
		Name:        username,
		IsApproved:  false,
		IsAdmin:     false,
		Credentials: []webauthn.Credential{},
	}
	user, err := w.userService.AddUser(user)
	if err != nil {
		return nil, err
	}
	return w.beginRegistration(user, username)
}

func (w WebAuthNService) startRegisterExistingUser(username string) (*protocol.CredentialCreation, error) {
	user, err := w.userService.GetUser(username)
	if err != nil {
		return nil, err
	}
	return w.beginRegistration(user, username)
}

func (w WebAuthNService) userExists(username string) bool {
	_, err := w.userService.GetUser(username)
	return err == nil
}

func (w WebAuthNService) beginRegistration(user *models.User, username string) (*protocol.CredentialCreation, error) {
	registerOpts := func(credOptions *protocol.PublicKeyCredentialCreationOptions) {
		credOptions.CredentialExcludeList = user.ExcludedCredentials()
	}

	options, sessionData, err := w.webAuthn.BeginRegistration(user, registerOpts)
	if err != nil {
		return nil, err
	}

	if err = w.sessionStore.storeSession(username, *sessionData); err != nil {
		return nil, err
	}

	return options, nil
}

func (w *WebAuthNService) finishRegistration(username string, request *http.Request) (*models.User, error) {
	sessionData, err := w.sessionStore.getSession(username)
	if err != nil {
		return nil, err
	}
	if time.Now().After(sessionData.Expires) {
		return nil, fmt.Errorf("Sessiondata not found or expired")
	}

	user, err := w.userService.GetUser(username)
	if err != nil {
		return nil, err
	}

	credential, err := w.webAuthn.FinishRegistration(user, sessionData, request)
	if err != nil {
		return nil, err
	}

	user.AddCredential(*credential)

	if !w.userService.HasApprovedAndAdminUser() {
		user.IsApproved = true
		user.IsAdmin = true
	}

	user, err = w.userService.SaveUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (w WebAuthNService) startLogin(username string, request *http.Request) (*protocol.CredentialAssertion, error) {
	user, err := w.userService.GetUser(username)
	if err != nil {
		return nil, err
	}

	options, sessionData, err := w.webAuthn.BeginLogin(user)
	if err != nil {
		return nil, err
	}

	if err = w.sessionStore.storeSession(username, *sessionData); err != nil {
		return nil, err
	}

	w.userService.SaveUser(user)
	return options, nil
}

func (w WebAuthNService) finishLogin(username string, request *http.Request) (*models.User, string, error) {
	sessionData, err := w.sessionStore.getSession(username)
	if err != nil {
		return nil, "", err
	}

	if time.Now().After(sessionData.Expires) {
		return nil, "", fmt.Errorf("Sessiondata not found or expired")
	}

	user, err := w.userService.GetUser(username)
	if err != nil {
		return nil, "", err
	}

	cred, err := protocol.ParseCredentialRequestResponseBody(request.Body)
	if err != nil {
		return nil, "", err
	}

	_, err = w.webAuthn.ValidateLogin(user, sessionData, cred)
	if err != nil {
		return nil, "", err
	}

	user.CredentialId = cred.ID

	user, err = w.userService.SaveUser(user)
	if err != nil {
		return nil, "", err
	}

	token := w.jwtService.GenerateToken(*user)
	return user, token, nil
}

func (w WebAuthNService) changePassword(ctx context.Context, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return w.userService.ChangePassword(ctx, username, string(hashedPassword))
}

func (w WebAuthNService) passwordLogin(ctx context.Context, username, password string) (*models.User, string, error) {
	if err := w.userService.CheckLogin(username, password); err != nil {
		return nil, "", err
	}
	user, err := w.userService.GetUser(username)
	if err != nil {
		return nil, "", err
	}

	token := w.jwtService.GenerateToken(*user)
	return user, token, nil
}

func (w WebAuthNService) startDiscoverableLogin(request *http.Request) (*protocol.CredentialAssertion, error) {
	options, sessionData, err := w.webAuthn.BeginDiscoverableLogin()
	if err != nil {
		return nil, err
	}

	sessionID := sessionData.Challenge
	if err = w.sessionStore.storeSession(sessionID, *sessionData); err != nil {
		return nil, err
	}

	return options, nil
}

func (w WebAuthNService) finishDiscoverableLogin(request *http.Request) (*models.User, string, error) {
	response, err := protocol.ParseCredentialRequestResponseBody(request.Body)
	if err != nil {
		return nil, "", err
	}

	sessionID := response.Response.CollectedClientData.Challenge
	sessionData, err := w.sessionStore.getSession(sessionID)
	if err != nil {
		return nil, "", err
	}

	if time.Now().After(sessionData.Expires) {
		return nil, "", fmt.Errorf("session expired")
	}

	credentialId := response.ID
	if len(credentialId) == 0 {
		return nil, "", fmt.Errorf("credentialId is required for discoverable login")
	}

	user, err := w.userService.GetUserByCredentialId(credentialId)
	if err != nil {
		return nil, "", fmt.Errorf("user not found: %w", err)
	}

	_, err = w.webAuthn.ValidateDiscoverableLogin(func(rawID, userHandle []byte) (webauthn.User, error) {
		return user, nil
	}, sessionData, response)
	if err != nil {
		return nil, "", err
	}

	user, err = w.userService.SaveUser(user)
	if err != nil {
		return nil, "", err
	}

	token := w.jwtService.GenerateToken(*user)
	return user, token, nil
}
