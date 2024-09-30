package security

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/kordondev/equipment-watchdog/models"
	"github.com/kordondev/equipment-watchdog/url"
)

type Service interface {
	startRegisterExistingUser(string) (*protocol.CredentialCreation, error)
	startRegisterNewUser(string) (*protocol.CredentialCreation, error)
	finishRegistration(string, *http.Request) (*models.User, error)
	startLogin(string, *http.Request) (*protocol.CredentialAssertion, error)
	finishLogin(string, *http.Request) (*models.User, string, error)
	passwordLogin(ctx context.Context, username, password string) (*models.User, string, error)
	changePassword(ctx context.Context, username, password string) error
}

type Controller struct {
	service Service
	domain  string
}

func NewController(baseUrl *gin.RouterGroup, service Service, domain string, authorizeMiddleware gin.HandlerFunc) error {

	ctrl := Controller{
		service: service,
		domain:  domain,
	}
	baseUrl.GET("/register/:username", ctrl.startRegisterNewUser)
	baseUrl.POST("/register/:username", ctrl.finishRegistration)
	baseUrl.GET("/login/:username", ctrl.startLogin)
	baseUrl.POST("/login/:username", ctrl.finishLogin)
	baseUrl.POST("/logout", ctrl.logout)
	baseUrl.POST("/password-login", ctrl.passwordLogin)

	baseUrl.PATCH("/change-password", authorizeMiddleware, ctrl.changePassword)
	baseUrl.GET("/add-authentication", authorizeMiddleware, ctrl.startRegisterExistingUser)
	baseUrl.POST("/add-authentication", authorizeMiddleware, ctrl.finishRegisterExistingUser)

	return nil
}

func (ctrl Controller) startRegisterNewUser(c *gin.Context) {
	username := c.Param("username")

	options, err := ctrl.service.startRegisterNewUser(username)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, options)
}

func (ctrl Controller) startRegisterExistingUser(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	options, err := ctrl.service.startRegisterExistingUser(username)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, options)
}

func (ctrl Controller) finishRegisterExistingUser(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := ctrl.service.finishRegistration(username, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
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

type changePasswordRequest struct {
	Password string `json:"password"`
}

func (ctrl Controller) changePassword(c *gin.Context) {
	var p changePasswordRequest
	if err := c.BindJSON(&p); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	username := c.GetString("username")
	if username == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := ctrl.service.changePassword(c, username, p.Password)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

type passwordLogin struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func (ctrl Controller) passwordLogin(c *gin.Context) {
	var login passwordLogin
	if err := c.BindJSON(&login); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, token, err := ctrl.service.passwordLogin(c, login.Username, login.Password)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
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
