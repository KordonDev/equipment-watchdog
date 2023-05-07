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
)

type Service interface {
	StartRegister(string) (*protocol.CredentialCreation, *webauthn.SessionData, error)
	FinishRegistration(string, webauthn.SessionData, *http.Request) (*models.User, error)
	StartLogin(string, *http.Request) (*protocol.CredentialAssertion, *webauthn.SessionData, error)
	FinishLogin(string, webauthn.SessionData, *http.Request, *gin.Context) (*models.User, error)
	Logout(*gin.Context)
}

type Controller struct {
	service      Service
	sessionStore *session.Store
}

func NewController(baseUrl *gin.RouterGroup, service Service) error {
	sessionStore, err := session.NewStore()
	if err != nil {
		log.Fatal("Error creating sessionStore", err)
		return fmt.Errorf("could not create session-store: %w", err)
	}

	ctrl := Controller{
		service:      service,
		sessionStore: sessionStore,
	}
	baseUrl.GET("/register/:username", ctrl.StartRegister)
	baseUrl.POST("/register/:username", ctrl.FinishRegistration)
	baseUrl.GET("/login/:username", ctrl.StartLogin)
	baseUrl.POST("/login/:username", ctrl.FinishLogin)
	baseUrl.POST("/logout", ctrl.Logout)

	return nil
}

func (ctrl Controller) StartRegister(c *gin.Context) {
	username := c.Param("username")

	options, sessionData, err := ctrl.service.StartRegister(username)
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

func (ctrl Controller) FinishRegistration(c *gin.Context) {
	username := c.Param("username")

	sessionData, err := ctrl.sessionStore.GetWebauthnSession("registration", c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := ctrl.service.FinishRegistration(username, sessionData, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl Controller) StartLogin(c *gin.Context) {
	username := c.Param("username")

	options, sessionData, err := ctrl.service.StartLogin(username, c.Request)

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

func (ctrl Controller) FinishLogin(c *gin.Context) {
	username := c.Param("username")

	sessionData, err := ctrl.sessionStore.GetWebauthnSession("authentication", c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := ctrl.service.FinishLogin(username, sessionData, c.Request, c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl Controller) Logout(c *gin.Context) {
	ctrl.service.Logout(c)
	c.JSON(http.StatusUnauthorized, gin.H{
		"redirect": "login",
	})
}
