package security

import (
	"fmt"
	"net/http"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/kordondev/equipment-watchdog/models"
)

type userService interface {
	GetUser(string) (*models.User, error)
	AddUser(*models.User) (*models.User, error)
	SaveUser(*models.User) (*models.User, error)
	HasApprovedAndAdminUser() bool
}

type WebAuthNService struct {
	webAuthn    *webauthn.WebAuthn
	jwtService  *JwtService
	userService userService
	domain      string
}

func NewWebAuthNService(userService userService, origin string, domain string, jwtService *JwtService) (*WebAuthNService, error) {
	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "equipment watchdog", // Display Name for your site
		RPID:          domain,               // Generally the domain name for your site
		RPOrigin:      origin,               // The origin URL for WebAuthn requests
	})

	if err != nil {
		return nil, fmt.Errorf("could not create webAuth: %w", err)
	}

	return &WebAuthNService{
		webAuthn:    webAuthn,
		jwtService:  jwtService,
		userService: userService,
		domain:      domain,
	}, nil
}

func (w WebAuthNService) startRegister(username string) (*protocol.CredentialCreation, *webauthn.SessionData, error) {
	user, err := w.userService.GetUser(username)
	if err != nil {
		user = &models.User{
			Name:        username,
			IsApproved:  false,
			IsAdmin:     false,
			Credentials: []webauthn.Credential{},
		}
		user, err = w.userService.AddUser(user)
		if err != nil {
			return nil, nil, err
		}
	}

	registerOpts := func(credOptions *protocol.PublicKeyCredentialCreationOptions) {
		credOptions.CredentialExcludeList = user.ExcludedCredentials()
	}

	options, sessionData, err := w.webAuthn.BeginRegistration(
		user,
		registerOpts,
	)
	if err != nil {
		return nil, nil, err
	}

	return options, sessionData, nil
}

func (w *WebAuthNService) finishRegistration(username string, sessionData webauthn.SessionData, request *http.Request) (*models.User, error) {

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

func (w WebAuthNService) startLogin(username string, request *http.Request) (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	user, err := w.userService.GetUser(username)
	if err != nil {
		return nil, nil, err
	}

	options, sessionData, err := w.webAuthn.BeginLogin(user)
	if err != nil {
		return nil, nil, err
	}

	w.userService.SaveUser(user)
	return options, sessionData, nil
}

func (w WebAuthNService) finishLogin(username string, sessionData webauthn.SessionData, request *http.Request) (*models.User, string, error) {
	user, err := w.userService.GetUser(username)
	if err != nil {
		return nil, "", err
	}

	_, err = w.webAuthn.FinishLogin(user, sessionData, request)
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
