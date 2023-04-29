package members

import (
	"net/http"

	"github.com/cloudflare/cfssl/log"
	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/equipment"
	"github.com/kordondev/equipment-watchdog/models"
	"github.com/kordondev/equipment-watchdog/url"
)

type MemberDatabase interface {
	GetMemberById(id uint64) (*models.Member, error)
	GetAllMember() ([]*models.Member, error)
	DeleteMember(*models.Member) error
	CreateMember(*models.Member) (*models.Member, error)
	SaveMember(*models.Member) error
}
type MemberService struct {
	db               MemberDatabase
	equipmentService *equipment.EquipmentService
}

func NewMemberService(database MemberDatabase, equipmentService *equipment.EquipmentService) *MemberService {
	return &MemberService{
		db:               database,
		equipmentService: equipmentService,
	}
}

func (s *MemberService) GetAllMembers(c *gin.Context) {
	m, err := s.db.GetAllMember()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, m)
}

func (s *MemberService) GetMemberById(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	m, err := s.db.GetMemberById(id)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, m)
}

func (s *MemberService) UpdateMember(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	em, err := s.db.GetMemberById(id)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var um models.Member
	if err := c.BindJSON(&um); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	eqIds := make([]uint64, 0)

	for _, eT := range models.GroupWithEquipment[um.Group] {
		if um.Equipment[eT] != nil {
			eqIds = append(eqIds, um.Equipment[eT].Id)
		}
	}
	equipments, err := s.equipmentService.GetAllByIds(eqIds)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	um.Id = em.Id
	um.Equipment = um.ListToMap(equipments, um.Id)
	err = s.db.SaveMember(&um)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, um)
}

func (s *MemberService) CreateMember(c *gin.Context) {
	var m models.Member
	if err := c.BindJSON(&m); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	cm, err := s.db.CreateMember(&m)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, cm)
}

func (s *MemberService) DeleteById(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	err = s.db.DeleteMember(&models.Member{Id: id})
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *MemberService) GetAllGroups(c *gin.Context) {
	c.JSON(http.StatusOK, models.GroupWithEquipment)
}
