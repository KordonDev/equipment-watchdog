package orders

import (
	"net/http"

	"github.com/cloudflare/cfssl/log"
	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/models"
	"github.com/kordondev/equipment-watchdog/url"
	"gorm.io/gorm"
)

type OrderService struct {
	db *orderDB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{
		db: newOrderDB(db),
	}
}

func (s OrderService) Create(c *gin.Context) {
	var o models.Order
	if err := c.BindJSON(&o); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	co, err := s.db.create(o)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, co)
}

func (s OrderService) GetById(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	o, err := s.db.getById(id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, o)
}

func (s OrderService) Update(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	eo, err := s.db.getById(id)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var uo models.Order
	if err := c.BindJSON(uo); err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	uo.ID = eo.ID
	uo.CreatedAt = eo.CreatedAt
	err = s.db.save(&uo)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, uo)
}

func (s OrderService) Delete(c *gin.Context) {
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
