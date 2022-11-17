package members

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/cloudflare/cfssl/log"
	"github.com/gin-gonic/gin"
)

type MemberService struct {
	memberDB *memberDB
}

func NewMemberService(memberDB *memberDB) *MemberService {
	return &MemberService{
		memberDB: memberDB,
	}
}

func (s *MemberService) GetAllMembers(c *gin.Context) {
	m, err := s.memberDB.GetAllMember()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, m)
}

func (s *MemberService) GetMemberById(c *gin.Context) {
	id := c.Param("id")
	m, err := s.getMember(id)

	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, m)
}

func (s *MemberService) UpdateMember(c *gin.Context) {
	id := c.Param("id")
	em, err := s.getMember(id)

	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var um member
	if err := c.BindJSON(&um); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	um.Id = em.Id
	err = s.memberDB.SaveMember(&um)
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

	cm, err := s.memberDB.CreateMember(&m)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, cm)
}

func (s *MemberService) DeleteById(c *gin.Context) {
	id := c.Param("id")

	idN, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	err = s.memberDB.db.Delete(&dbMember{}, idN).Error

	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.Status(http.StatusOK)

}

func (s *MemberService) getMember(id string) (*member, error) {
	idN, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		return nil, errors.New("id as number could not be found")
	}

	return s.memberDB.GetMemberById(idN)
}
