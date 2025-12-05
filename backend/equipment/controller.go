package equipment

import (
	"net/http"

	"github.com/cloudflare/cfssl/log"
	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/models"
	"github.com/kordondev/equipment-watchdog/url"
)

type Service interface {
	getEquipmentById(uint64) (*models.Equipment, error)
	getAllEquipment() ([]*models.Equipment, error)
	getFreeEquipment() (map[models.EquipmentType][]*models.Equipment, error)
	getAllEquipmentByType(string) ([]*models.Equipment, error)
}

type Controller struct {
	service      Service
	changeWriter ChangeWriter
}

type ChangeWriter interface {
	Save(models.Change, *gin.Context) (*models.Change, error)
}

func NewController(baseRoute *gin.RouterGroup, service Service, changeWriter ChangeWriter) {

	ctrl := Controller{
		service:      service,
		changeWriter: changeWriter,
	}

	equipmentRoute := baseRoute.Group("/equipment")
	{
		equipmentRoute.GET("/:id", ctrl.getEquipmentById)
		equipmentRoute.GET("/", ctrl.getAllEquipment)
		equipmentRoute.GET("/type/:type", ctrl.getAllEquipmentByType)
		equipmentRoute.GET("/free", ctrl.getFreeEquipment)
	}

}

func (ctrl Controller) getEquipmentById(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	e, err := ctrl.service.getEquipmentById(id)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, e)
}

func (ctrl Controller) getAllEquipment(c *gin.Context) {
	e, err := ctrl.service.getAllEquipment()
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, e)
}

func (ctrl Controller) getFreeEquipment(c *gin.Context) {
	equipments, err := ctrl.service.getFreeEquipment()
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, equipments)
}

func (ctrl Controller) getAllEquipmentByType(c *gin.Context) {
	eType := c.Param("type")
	if eType == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	es, err := ctrl.service.getAllEquipmentByType(eType)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, es)
}
