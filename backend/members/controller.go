package members

import (
	"github.com/cloudflare/cfssl/log"
	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/models"
	"github.com/kordondev/equipment-watchdog/url"
	"net/http"
)

type Service interface {
	createMember(*models.Member) (*models.Member, error)
	getMemberById(uint64) (*models.Member, error)
	updateMember(uint64, *models.Member) ([]uint64, error)
	deleteMemberById(uint64) error
	getAllGroups() map[models.Group][]models.EquipmentType
	getAllMembers() ([]*models.Member, error)
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

	changedEquipments, err := ctrl.service.updateMember(id, &um)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	for _, id := range changedEquipments {
		ctrl.changeWriter.Save(models.Change{
			Action:      models.UpdateMember,
			MemberId:    um.Id,
			EquipmentId: id,
		}, c)
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
