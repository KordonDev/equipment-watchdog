package equipment

import (
	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
	"net/http"

	"github.com/cloudflare/cfssl/log"
	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/url"
)

type EquipmentService struct {
	db *equipmentDB
}

func NewEquipmentService(db *gorm.DB) *EquipmentService {
	return &EquipmentService{
		db: newEquipmentDB(db),
	}
}

func (s *EquipmentService) GetEquipmentById(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	e, err := s.db.getById(id)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, e)
}

func (s *EquipmentService) GetAllEquipmentByType(c *gin.Context) {
	eType := c.Param("type")
	if eType == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	es, err := s.db.getByType(eType)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, es)
}

func (s *EquipmentService) CreateEquipment(c *gin.Context) {
	var e models.Equipment
	if err := c.BindJSON(&e); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ce, err := s.db.Create(&e)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, ce)
}

func (s *EquipmentService) DeleteEquipment(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = s.db.delete(id)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
