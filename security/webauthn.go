package security

import (
	"fmt"
	"net/http"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/models"
)

type WebAuthNService struct {
	webAuthn   *webauthn.WebAuthn
	jwtService *JwtService
	userDB     *userDB
	domain     string
}

const AUTHORIZATION_COOKIE_KEY = "Authorization"

func NewWebAuthNService(userDB *userDB, origin string, domain string, jwtService *JwtService) (*WebAuthNService, error) {
	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "equipment watchdog", // Display Name for your site
		RPID:          domain,               // Generally the domain name for your site
		RPOrigin:      origin,               // The origin URL for WebAuthn requests
	})

	if err != nil {
		return nil, fmt.Errorf("could not create webAuth: %w", err)
	}

	return &WebAuthNService{
		webAuthn:   webAuthn,
		jwtService: jwtService,
		userDB:     userDB,
		domain:     domain,
	}, nil
}

// FIXME: not really nice solution to cancle on this way the context - channel issue
// try to make clean architecture
func (w WebAuthNService) StartRegister(username string) (*protocol.CredentialCreation, *webauthn.SessionData, error) {
	user, err := w.userDB.GetUser(username)
	if err != nil {
		user = &models.User{
			Name:        username,
			IsApproved:  false,
			IsAdmin:     false,
			Credentials: []webauthn.Credential{},
		}
		user, err = w.userDB.AddUser(user)
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

func (w *WebAuthNService) FinishRegistration(username string, sessionData webauthn.SessionData, request *http.Request) (*models.User, error) {

	user, err := w.userDB.GetUser(username)
	if err != nil {
		return nil, err
	}

	credential, err := w.webAuthn.FinishRegistration(user, sessionData, request)
	if err != nil {
		return nil, err
	}

	user.AddCredential(*credential)

	if !w.userDB.HasApprovedAndAdminUser() {
		user.IsApproved = true
		user.IsAdmin = true
	}

	user, err = w.userDB.SaveUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (w WebAuthNService) StartLogin(username string, request *http.Request) (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	user, err := w.userDB.GetUser(username)
	if err != nil {
		return nil, nil, err
	}

	options, sessionData, err := w.webAuthn.BeginLogin(user)
	if err != nil {
		return nil, nil, err
	}

	w.userDB.SaveUser(user)
	return options, sessionData, nil
}

func (w WebAuthNService) FinishLogin(username string, sessionData webauthn.SessionData, request *http.Request, c *gin.Context) (*models.User, error) {
	user, err := w.userDB.GetUser(username)
	if err != nil {
		return nil, err
	}

	_, err = w.webAuthn.FinishLogin(user, sessionData, request)
	if err != nil {
		return nil, err
	}

	user, err = w.userDB.SaveUser(user)
	if err != nil {
		return nil, err
	}

	token := w.jwtService.GenerateToken(*user)
	w.jwtService.SetCookie(c, token)
	return user, nil
}

func (w WebAuthNService) Logout(c *gin.Context) {
	c.SetCookie(AUTHORIZATION_COOKIE_KEY, "", -1, "/", w.domain, true, true)
}
