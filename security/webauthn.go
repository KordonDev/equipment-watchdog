package security

import (
	"log"
	"net/http"

	"github.com/duo-labs/webauthn.io/session"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gin-gonic/gin"
)

type UserRequest struct {
	username string `json:"username"`
}

type WebAuthNService struct {
	webAuthn     *webauthn.WebAuthn
	sessionStore *session.Store
	jwtService   JWTService
	userDB       *userDB
}

func NewWebAuthNService(userDB *userDB) *WebAuthNService {
	var err error
	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "equipment watchdog",    // Display Name for your site
		RPID:          "localhost",             // Generally the domain name for your site
		RPOrigin:      "http://localhost:8080", // The origin URL for WebAuthn requests
	})

	if err != nil {
		log.Fatal("Error creating webAuthn", err)
	}

	sessionStore, err := session.NewStore()
	if err != nil {
		log.Fatal("Error creating sessionStore", err)
	}

	jwtService := JWTAuthService()

	return &WebAuthNService{webAuthn: webAuthn, sessionStore: sessionStore, jwtService: jwtService, userDB: userDB}
}

func (w *WebAuthNService) StartRegister(c *gin.Context) {
	username := c.Param("username")

	user, err := w.userDB.GetUser(username)
	if err != nil {
		user = &User{
			name:        username,
			credentials: []webauthn.Credential{},
		}
		w.userDB.AddUser(user)
	}

	// TODO: go PublicKeyCredentialCreationOptions
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

	c.SetCookie("Authorization2", token, 60*100, "/", "/", true, true)
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
