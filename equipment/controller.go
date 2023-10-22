package equipment

import (
	"github.com/cloudflare/cfssl/log"
	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/models"
	"github.com/kordondev/equipment-watchdog/url"
	"net/http"
)

type Service interface {
	createEquipment(models.Equipment) (*models.Equipment, error)
	getEquipmentById(uint64) (*models.Equipment, error)
	deleteEquipment(uint64) error
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
		equipmentRoute.GET("/type/:type", ctrl.getAllEquipmentByType)
		equipmentRoute.GET("/free", ctrl.getFreeEquipment)
		equipmentRoute.POST("/", ctrl.createEquipment)
		equipmentRoute.DELETE("/:id", ctrl.deleteEquipment)
	}

}

func (ctrl Controller) createEquipment(c *gin.Context) {
	var e models.Equipment
	if err := c.BindJSON(&e); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ce, err := ctrl.service.createEquipment(e)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctrl.changeWriter.Save(models.Change{
		Action:      models.CreateEquipment,
		EquipmentId: ce.Id,
		MemberId:    ce.MemberID,
	}, c)

	c.JSON(http.StatusCreated, ce)
}

func (ctrl Controller) getEquipmentById(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		log.Error(err)
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

func (ctrl Controller) deleteEquipment(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = ctrl.service.deleteEquipment(id)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctrl.changeWriter.Save(models.Change{
		Action:      models.DeleteEquipment,
		EquipmentId: id,
	}, c)

	c.Status(http.StatusOK)
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
