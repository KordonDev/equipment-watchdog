package members

import (
	"net/http"

	"github.com/cloudflare/cfssl/log"
	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/url"
)

type MemberService struct {
	db *memberDB
}

func NewMemberService(memberDB *memberDB) *MemberService {
	return &MemberService{
		db: memberDB,
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

	var um member
	if err := c.BindJSON(&um); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	um.Id = em.Id
	err = s.db.SaveMember(&um)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, um)
}

func (s *MemberService) CreateMember(c *gin.Context) {
	var m member
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

	err = s.db.db.Delete(&dbMember{}, id).Error
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *MemberService) GetAllGroups(c *gin.Context) {
	c.JSON(http.StatusOK, GroupWithEquipment)
}
