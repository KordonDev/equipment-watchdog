package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/models"
	"github.com/kordondev/equipment-watchdog/security"
)

type Service interface {
	GetAll() ([]*models.User, error)
	GetUserWithToken(string, *gin.Context) (*models.User, error)
	ToggleApprove(string) (*models.User, error)
	ToggleAdmin(string) (*models.User, error)
}

type Controller struct {
	service Service
}

// TODO: or security as parameter?
func NewController(baseUrl *gin.RouterGroup, service Service) error {

	ctrl := Controller{
		service: service,
	}

	baseUrl.GET("/me", ctrl.GetMe)
	userRoute := baseUrl.Group("/users")
	{
		userRoute.GET("/", security.AdminOnlyMiddleware(), ctrl.GetAll)

		userRoute.PATCH("/:username/toggle-approve", security.AdminOnlyMiddleware(), ctrl.ToggleApprove)
		userRoute.PATCH("/:username/toggle-admin", security.AdminOnlyMiddleware(), ctrl.ToggleAdmin)
	}
	return nil
}

func (ctrl Controller) GetMe(c *gin.Context) {
	username := c.GetString("username")
	user, err := ctrl.service.GetUserWithToken(username, c)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl Controller) GetAll(c *gin.Context) {
	users, err := ctrl.service.GetAll()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

func (ctrl Controller) ToggleApprove(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := ctrl.service.ToggleApprove(username)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl Controller) ToggleAdmin(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := ctrl.service.ToggleAdmin(username)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, user)
}
