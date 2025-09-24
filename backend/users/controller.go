package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/models"
	"github.com/kordondev/equipment-watchdog/security"
	"github.com/kordondev/equipment-watchdog/url"
)

type Service interface {
	GetUser(string) (*models.User, error)
	getAll() ([]*models.User, error)
	getTokenForUser(*models.User) string
	toggleApprove(string) (*models.User, error)
	toggleAdmin(string) (*models.User, error)
}

type Controller struct {
	service Service
	domain  string
}

func NewController(baseUrl *gin.RouterGroup, service Service, domain string) {

	ctrl := Controller{
		service: service,
		domain:  domain,
	}

	baseUrl.GET("/me", ctrl.getMe)
	userRoute := baseUrl.Group("/users")
	{
		userRoute.GET("/", security.AdminOnlyMiddleware(), ctrl.getAll)

		userRoute.PATCH("/:username/toggle-approve", security.AdminOnlyMiddleware(), ctrl.toggleApprove)
		userRoute.PATCH("/:username/toggle-admin", security.AdminOnlyMiddleware(), ctrl.toggleAdmin)
	}
}

func (ctrl Controller) getMe(c *gin.Context) {
	username := c.GetString("username")

	user, err := ctrl.service.GetUser(username)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	token := ctrl.service.getTokenForUser(user)
	url.SetCookie(c, token, ctrl.domain)
	c.JSON(http.StatusOK, user)
}

func (ctrl Controller) getAll(c *gin.Context) {
	users, err := ctrl.service.getAll()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

func (ctrl Controller) toggleApprove(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := ctrl.service.toggleApprove(username)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl Controller) toggleAdmin(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := ctrl.service.toggleAdmin(username)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, user)
}
