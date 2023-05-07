package equipment

import (
	"github.com/cloudflare/cfssl/log"
	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/models"
	"github.com/kordondev/equipment-watchdog/url"
	"net/http"
)

type Service interface {
	CreateEquipment(models.Equipment) (*models.Equipment, error)
	GetEquipmentById(uint64) (*models.Equipment, error)
	DeleteEquipment(uint64) error
  GetFreeEquipment() (map[models.EquipmentType][]*models.Equipment, error)
  GetAllEquipmentByType(string)([]*models.Equipment, error)
}

type Controller struct {
	service Service
}

func NewController(baseRoute *gin.RouterGroup, service Service) {

	ctrl := Controller{
		service: service,
	}

	equipmentRoute := baseRoute.Group("/equipment")
	{
		equipmentRoute.GET("/:id", ctrl.GetEquipmentById)
		equipmentRoute.GET("/type/:type", ctrl.GetAllEquipmentByType)
		equipmentRoute.GET("/free", ctrl.GetFreeEquipment)
		equipmentRoute.POST("/", ctrl.CreateEquipment)
		equipmentRoute.DELETE("/:id", ctrl.DeleteEquipment)
	}

}

func (ctrl Controller) CreateEquipment(c *gin.Context) {
	var e models.Equipment
	if err := c.BindJSON(&e); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ce, err := ctrl.service.CreateEquipment(e)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, ce)
}

func (ctrl Controller) GetEquipmentById(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	e, err := ctrl.service.GetEquipmentById(id)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, e)
}

func (ctrl Controller) DeleteEquipment(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = ctrl.service.DeleteEquipment(id)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (ctrl Controller) GetFreeEquipment(c *gin.Context) {
  equipments, err := ctrl.service.GetFreeEquipment()
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, equipments)
}

func (ctrl Controller) GetAllEquipmentByType(c *gin.Context) {
	eType := c.Param("type")
	if eType == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	es, err := ctrl.service.GetAllEquipmentByType(eType)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, es)

}
