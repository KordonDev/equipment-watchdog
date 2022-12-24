package equipment

import (
	"net/http"

	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"

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

func (s *EquipmentService) GetAllByIds(ids []uint64) ([]*models.Equipment, error) {
	return s.db.getAllByIds(ids)
}

func (s *EquipmentService) FreeEquipment(c *gin.Context) {
	equipment, err := s.db.getFreeEquipment()

	equipments := make(map[models.EquipmentType][]*models.Equipment)
	for _, e := range equipment {
		equipments[e.Type] = append(equipments[e.Type], e)
	}
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, equipments)
}
