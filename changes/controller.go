package changes

import (
	"net/http"

	"github.com/cloudflare/cfssl/log"
	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/models"
)

type Service interface {
	getAll() ([]*models.Change, error)
}

type Controller struct {
	service Service
}

func NewController(baseRoute *gin.RouterGroup, service Service) {

	ctrl := Controller{
		service: service,
	}

	changesRoute := baseRoute.Group("/changes")
	{
		changesRoute.GET("/", ctrl.getAllChanges)
	}

}

func (ctrl Controller) getAllChanges(c *gin.Context) {
	equipments, err := ctrl.service.getAll()
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, equipments)

}
