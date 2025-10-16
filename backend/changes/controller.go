package changes

import (
	"net/http"

	"github.com/cloudflare/cfssl/log"
	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/models"
	"github.com/kordondev/equipment-watchdog/url"
)

type Service interface {
	getAll() ([]string, error)
	getForEquipment(uint64) ([]string, error)
	getForOrder(uint64) ([]string, error)
	getForMember(uint64) ([]string, error)
	getRecent() ([]*models.Change, error)
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
		changesRoute.GET("/recent", ctrl.getRecentChanges)
		changesRoute.GET("/members/:id", ctrl.getForMember)
		changesRoute.GET("/orders/:id", ctrl.getForOrder)
		changesRoute.GET("/equipments/:id", ctrl.getForEquipment)
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

func (ctrl Controller) getForEquipment(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	cs, err := ctrl.service.getForEquipment(id)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, cs)
}

func (ctrl Controller) getForOrder(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	cs, err := ctrl.service.getForOrder(id)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, cs)
}

func (ctrl Controller) getForMember(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	cs, err := ctrl.service.getForMember(id)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, cs)
}

func (ctrl Controller) getRecentChanges(c *gin.Context) {
	cs, err := ctrl.service.getRecent()
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, cs)
}
