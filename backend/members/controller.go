package members

import (
	"net/http"

	"github.com/cloudflare/cfssl/log"
	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/models"
	"github.com/kordondev/equipment-watchdog/url"
)

type Service interface {
	createMember(*models.Member) (*models.Member, error)
	getMemberById(uint64) (*models.Member, error)
	updateMember(uint64, *models.Member) error
	deleteMemberById(uint64) error
	getAllGroups() map[models.Group][]models.EquipmentType
	getAllMembers() ([]*models.Member, error)
	saveEquipmentForMember(uint64, models.EquipmentType, models.Equipment) (*models.Equipment, *models.Equipment, error)
	removeEquipmentFromMember(uint64, models.EquipmentType) (*models.Equipment, error)
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

	membersRoute := baseRoute.Group("/members")
	{
		membersRoute.GET("/", ctrl.getAllMembers)
		membersRoute.GET("/:id", ctrl.getMemberById)
		membersRoute.POST("/", ctrl.createMember)
		membersRoute.PUT("/:id", ctrl.updateMember)
		membersRoute.DELETE("/:id", ctrl.deleteMemberById)
		membersRoute.GET("/groups", ctrl.getAllGroups)
		membersRoute.POST("/:id/:equipmentType", ctrl.saveEquipmentForMember)
		membersRoute.DELETE("/:id/:equipmentType", ctrl.removeEquipmentFromMember)
	}
}

func (ctrl Controller) getAllMembers(c *gin.Context) {
	m, err := ctrl.service.getAllMembers()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, m)

}

func (ctrl Controller) createMember(c *gin.Context) {
	var m models.Member
	if err := c.BindJSON(&m); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	cm, err := ctrl.service.createMember(&m)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctrl.changeWriter.Save(models.Change{
		Action:   models.CreateMember,
		MemberId: cm.Id,
	}, c)

	c.JSON(http.StatusCreated, cm)
}

func (ctrl Controller) getMemberById(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	m, err := ctrl.service.getMemberById(id)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, m)
}

func (ctrl Controller) updateMember(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if _, err = ctrl.service.getMemberById(id); err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var um models.Member
	if err := c.BindJSON(&um); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = ctrl.service.updateMember(id, &um)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, um)
}

func (ctrl Controller) deleteMemberById(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if err = ctrl.service.deleteMemberById(id); err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	ctrl.changeWriter.Save(models.Change{
		Action:   models.DeleteMember,
		MemberId: id,
	}, c)

	c.Status(http.StatusOK)
}

func (ctrl Controller) getAllGroups(c *gin.Context) {
	c.JSON(http.StatusOK, ctrl.service.getAllGroups())
}

func (ctrl Controller) saveEquipmentForMember(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	equipmentType := c.Param("equipmentType")
	if equipmentType == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var equipment models.Equipment
	if err := c.BindJSON(&equipment); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	savedEquipment, oldEquip, err := ctrl.service.saveEquipmentForMember(id, models.EquipmentType(equipmentType), equipment)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctrl.changeWriter.Save(models.Change{
		Action:         models.UpdateEquipmentOnMember,
		MemberId:       id,
		EquipmentId:    savedEquipment.Id,
		OldEquipmentId: oldEquip.Id,
	}, c)

	c.JSON(http.StatusOK, savedEquipment)
}

func (ctrl Controller) removeEquipmentFromMember(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	equipmentType := c.Param("equipmentType")
	if equipmentType == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	oldEquip, err := ctrl.service.removeEquipmentFromMember(id, models.EquipmentType(equipmentType))
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var oldEquipId uint64 = 0
	if oldEquip != nil {
		oldEquipId = oldEquip.Id
	}

	ctrl.changeWriter.Save(models.Change{
		Action:         models.UpdateEquipmentOnMember,
		MemberId:       id,
		OldEquipmentId: oldEquipId,
	}, c)

	c.JSON(http.StatusOK, oldEquip)
}
