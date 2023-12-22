package security

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/kordondev/equipment-watchdog/models"
	"github.com/kordondev/equipment-watchdog/url"
)

type Service interface {
	startRegister(string) (*protocol.CredentialCreation, error)
	finishRegistration(string, *http.Request) (*models.User, error)
	startLogin(string, *http.Request) (*protocol.CredentialAssertion, error)
	finishLogin(string, *http.Request) (*models.User, string, error)
}

type Controller struct {
	service Service
	domain  string
}

func NewController(baseUrl *gin.RouterGroup, service Service, domain string) error {

	ctrl := Controller{
		service: service,
		domain:  domain,
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

	options, err := ctrl.service.startRegister(username)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	/*err = ctrl.sessionStore.SaveWebauthnSession("registration", sessionData, c.Request, c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}*/

	c.JSON(http.StatusOK, options)
}

func (ctrl Controller) finishRegistration(c *gin.Context) {
	username := c.Param("username")

	user, err := ctrl.service.finishRegistration(username, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl Controller) startLogin(c *gin.Context) {
	username := c.Param("username")

	options, err := ctrl.service.startLogin(username, c.Request)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, options)
}

func (ctrl Controller) finishLogin(c *gin.Context) {
	username := c.Param("username")

	user, token, err := ctrl.service.finishLogin(username, c.Request)
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
