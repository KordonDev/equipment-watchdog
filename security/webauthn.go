package security

import (
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
	jwtService   JWTService
	userDB       *userDB
	domain       string
}

func NewWebAuthNService(userDB *userDB, origin string, domain string) *WebAuthNService {
	var err error
	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "equipment watchdog", // Display Name for your site
		RPID:          domain,               // Generally the domain name for your site
		RPOrigin:      origin,               // The origin URL for WebAuthn requests
	})

	if err != nil {
		log.Fatal("Error creating webAuthn", err)
	}

	sessionStore, err := session.NewStore()
	if err != nil {
		log.Fatal("Error creating sessionStore", err)
	}

	jwtService := JWTAuthService()

	return &WebAuthNService{webAuthn: webAuthn, sessionStore: sessionStore, jwtService: jwtService, userDB: userDB, domain: domain}
}

func (w *WebAuthNService) StartRegister(c *gin.Context) {
	username := c.Param("username")

	user, err := w.userDB.GetUser(username)
	if err != nil {
		user = &User{
			Name:        username,
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

func (w WebAuthNService) FinishRegistration(c *gin.Context) {
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
	w.userDB.SaveUser(user)
	c.Status(http.StatusOK)
}

func (w *WebAuthNService) StartLogin(c *gin.Context) {
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

func (w *WebAuthNService) FinishLogin(c *gin.Context) {
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

	token := w.jwtService.GenerateToken(username, true)
	w.userDB.SaveUser(user)

	c.SetCookie("Authorization", token, 60*100, "/", w.domain, true, true)
	c.Status(http.StatusOK)
}

func (w *WebAuthNService) Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", w.domain, true, true)
	c.JSON(http.StatusOK, gin.H{
		"redirect": "/index/Lisa",
	})
}
