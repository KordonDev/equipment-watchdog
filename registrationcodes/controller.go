package registrationcodes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/models"
)

type RCService interface {
	getRegistrationCode() (models.RegistrationCode, error)
}

type RegistrationCodesController struct {
	service RCService
}

func NewController(baseRoute *gin.RouterGroup, service RCService) {
	ctrl := RegistrationCodesController{
		service,
	}

	registrationCodesRoute := baseRoute.Group("/registrationCode")
	{
		registrationCodesRoute.GET("/", ctrl.getRegistrationCode)
	}
}

func (ctrl RegistrationCodesController) getRegistrationCode(c *gin.Context) {
	registrationCode, err := ctrl.service.getRegistrationCode()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, registrationCode)
}
