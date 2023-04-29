package security

import (
	"fmt"
	"log"
	"net/http"

	"github.com/duo-labs/webauthn.io/session"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gin-gonic/gin"
)

type WebAuthNService struct {
	webAuthn     *webauthn.WebAuthn
	sessionStore *session.Store
	jwtService   *JwtService
	userDB       *userDB
	domain       string
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

	sessionStore, err := session.NewStore()
	if err != nil {
		log.Fatal("Error creating sessionStore", err)
		return nil, fmt.Errorf("could not create session-store: %w", err)
	}

	return &WebAuthNService{
		webAuthn:     webAuthn,
		sessionStore: sessionStore,
		jwtService:   jwtService,
		userDB:       userDB,
		domain:       domain,
	}, nil
}

// FIXME: not really nice solution to cancle on this way the context - channel issue
// try to make clean architecture
func (w WebAuthNService) StartRegister(c *gin.Context) {
	username := c.Param("username")

	user, err := w.userDB.GetUser(username)
	if err != nil {
		user = &User{
			Name:        username,
			IsApproved:  false,
			IsAdmin:     false,
			Credentials: []webauthn.Credential{},
		}
		user, err = w.userDB.AddUser(user)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
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
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// store options
	err = w.sessionStore.SaveWebauthnSession("registration", sessionData, c.Request, c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, options)
}

func (w *WebAuthNService) FinishRegistration(c *gin.Context) {
	username := c.Param("username")

	user, err := w.userDB.GetUser(username)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	sessionData, err := w.sessionStore.GetWebauthnSession("registration", c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	credential, err := w.webAuthn.FinishRegistration(user, sessionData, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user.AddCredential(*credential)

	if !w.userDB.HasApprovedAndAdminUser() {
		user.IsApproved = true
		user.IsAdmin = true
	}

	user, err = w.userDB.SaveUser(user)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (w WebAuthNService) StartLogin(c *gin.Context) {
	username := c.Param("username")

	user, err := w.userDB.GetUser(username)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	options, sessionData, err := w.webAuthn.BeginLogin(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = w.sessionStore.SaveWebauthnSession("authentication", sessionData, c.Request, c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	w.userDB.SaveUser(user)

	c.JSON(http.StatusOK, options)
}

func (w WebAuthNService) FinishLogin(c *gin.Context) {
	username := c.Param("username")

	user, err := w.userDB.GetUser(username)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	sessionData, err := w.sessionStore.GetWebauthnSession("authentication", c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	_, err = w.webAuthn.FinishLogin(user, sessionData, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err = w.userDB.SaveUser(user)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	token := w.jwtService.GenerateToken(*user)

	w.jwtService.SetCookie(c, token)
	c.JSON(http.StatusOK, user)
}

func (w WebAuthNService) Logout(c *gin.Context) {
	c.SetCookie(AUTHORIZATION_COOKIE_KEY, "", -1, "/", w.domain, true, true)
	c.JSON(http.StatusUnauthorized, gin.H{
		"redirect": "login",
	})
}
