package registrationcodes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/models"
)

type Service interface {
  getRegistrationCode() models.RegistrationCode
}

type RegistrationCodesController struct {
  service Service
}

func NewController(baseRoute *gin.RouterGroup, service Service) {
  ctrl := RegistrationCodesController{
    service,
  }

	registrationCodesRoute := baseRoute.Group("/orders")
	{
		registrationCodesRoute.GET("/", ctrl.getRegistrationCode)
  }
}

func (ctrl RegistrationCodesController) getRegistrationCode(c *gin.Context) {
  registrationCode := ctrl.service.getRegistrationCode()

	c.JSON(http.StatusOK, registrationCode)
}
