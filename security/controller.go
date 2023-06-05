package security

import (
	"fmt"
	"log"
	"net/http"

	"github.com/duo-labs/webauthn.io/session"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/models"
	"github.com/kordondev/equipment-watchdog/url"
)

type Service interface {
	startRegister(string) (*protocol.CredentialCreation, *webauthn.SessionData, error)
	finishRegistration(string, webauthn.SessionData, *http.Request) (*models.User, error)
	startLogin(string, *http.Request) (*protocol.CredentialAssertion, *webauthn.SessionData, error)
	finishLogin(string, webauthn.SessionData, *http.Request) (*models.User, string, error)
}

type Controller struct {
	service      Service
	sessionStore *session.Store
	domain       string
}

func NewController(baseUrl *gin.RouterGroup, service Service, domain string) error {
	sessionStore, err := session.NewStore()
	if err != nil {
		log.Fatal("Error creating sessionStore", err)
		return fmt.Errorf("could not create session-store: %w", err)
	}

	ctrl := Controller{
		service:      service,
		sessionStore: sessionStore,
		domain:       domain,
	}
	baseUrl.GET("/register/:username", ctrl.startRegister)
	baseUrl.POST("/register/:username", ctrl.finishRegistration)
	baseUrl.GET("/login/:username", ctrl.startLogin)
	baseUrl.POST("/login/:username", ctrl.finishLogin)
	baseUrl.POST("/logout", ctrl.logout)

	return nil
}

func (ctrl Controller) startRegister(c *gin.Context) {
	username := c.Param("username")

	options, sessionData, err := ctrl.service.startRegister(username)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = ctrl.sessionStore.SaveWebauthnSession("registration", sessionData, c.Request, c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, options)
}

func (ctrl Controller) finishRegistration(c *gin.Context) {
	username := c.Param("username")

	sessionData, err := ctrl.sessionStore.GetWebauthnSession("registration", c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := ctrl.service.finishRegistration(username, sessionData, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl Controller) startLogin(c *gin.Context) {
	username := c.Param("username")

	options, sessionData, err := ctrl.service.startLogin(username, c.Request)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = ctrl.sessionStore.SaveWebauthnSession("authentication", sessionData, c.Request, c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, options)
}

func (ctrl Controller) finishLogin(c *gin.Context) {
	username := c.Param("username")

	sessionData, err := ctrl.sessionStore.GetWebauthnSession("authentication", c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, token, err := ctrl.service.finishLogin(username, sessionData, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	url.SetCookie(c, token, ctrl.domain)
	c.JSON(http.StatusOK, user)
}

func (ctrl Controller) logout(c *gin.Context) {
	url.RemoveCookie(c, ctrl.domain)
	c.JSON(http.StatusUnauthorized, gin.H{
		"redirect": "login",
	})
}
