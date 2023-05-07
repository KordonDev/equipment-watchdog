package members

import (
	"github.com/cloudflare/cfssl/log"
	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/models"
	"github.com/kordondev/equipment-watchdog/url"
	"net/http"
)

type Service interface {
	CreateMember(*models.Member) (*models.Member, error)
	GetMemberById(uint64) (*models.Member, error)
	UpdateMember(uint64, *models.Member) error
	DeleteMemberById(uint64) error
	GetAllGroups() map[models.Group][]models.EquipmentType
	GetAllMembers() ([]*models.Member, error)
}

type Controller struct {
	service Service
}

func NewController(baseRoute *gin.RouterGroup, service Service) {

	ctrl := Controller{
		service: service,
	}

	membersRoute := baseRoute.Group("/members")
	{
		membersRoute.GET("/", ctrl.GetAllMembers)
		membersRoute.GET("/:id", ctrl.GetMemberById)
		membersRoute.POST("/", ctrl.CreateMember)
		membersRoute.PUT("/:id", ctrl.UpdateMember)
		membersRoute.DELETE("/:id", ctrl.DeleteMemberById)
		membersRoute.GET("/groups", ctrl.GetAllGroups)
	}
}

func (ctrl Controller) GetAllMembers(c *gin.Context) {
	m, err := ctrl.service.GetAllMembers()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, m)

}

func (ctrl Controller) CreateMember(c *gin.Context) {
	var m models.Member
	if err := c.BindJSON(&m); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	cm, err := ctrl.service.CreateMember(&m)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, cm)
}

func (ctrl Controller) GetMemberById(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	m, err := ctrl.service.GetMemberById(id)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, m)
}

func (ctrl Controller) UpdateMember(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if _, err = ctrl.service.GetMemberById(id); err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var um models.Member
	if err := c.BindJSON(&um); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = ctrl.service.UpdateMember(id, &um)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, um)
}

func (ctrl Controller) DeleteMemberById(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if err = ctrl.service.DeleteMemberById(id); err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.Status(http.StatusOK)
}

func (ctrl Controller) GetAllGroups(c *gin.Context) {
	c.JSON(http.StatusOK, ctrl.service.GetAllGroups())
}
