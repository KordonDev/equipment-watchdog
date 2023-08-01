package orders

import (
	"github.com/cloudflare/cfssl/log"
	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/models"
	"github.com/kordondev/equipment-watchdog/url"
	"net/http"
)

type Service interface {
	getById(uint64) (models.Order, error)
	create(models.Order) (models.Order, error)
	update(uint64, models.Order) (models.Order, error)
	delete(uint64) error
	getAll(bool) ([]models.Order, error)
	fulfill(models.Order, string) (*models.Equipment, error)
}

type Controller struct {
	service Service
}

func NewController(baseRoute *gin.RouterGroup, service Service) {

	ctrl := Controller{
		service: service,
	}

	ordersRoute := baseRoute.Group("/orders")
	{
		ordersRoute.GET("/", ctrl.getAllNotFulfilled)
		ordersRoute.GET("/fullfilled", ctrl.getAllFulfilled)
		ordersRoute.GET("/:id", ctrl.getById)
		ordersRoute.POST("/", ctrl.create)
		ordersRoute.POST("/:registrationCode/toEquipment", ctrl.createEquipmentFromOrder)
		ordersRoute.PUT("/:id", ctrl.update)
		ordersRoute.DELETE("/:id", ctrl.delete)
	}
}

func (ctrl Controller) getAllNotFulfilled(c *gin.Context) {
	orders, err := ctrl.service.getAll(false)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (ctrl Controller) getAllFulfilled(c *gin.Context) {
	orders, err := ctrl.service.getAll(true)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (ctrl Controller) getById(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	o, err := ctrl.service.getById(id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, o)
}

func (ctrl Controller) create(c *gin.Context) {
	var o models.Order
	if err := c.BindJSON(&o); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	co, err := ctrl.service.create(o)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, co)
}

func (ctrl Controller) update(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
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

	order, err := ctrl.service.update(id, uo)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, order)
}

func (ctrl Controller) delete(c *gin.Context) {
	id, err := url.ParseToInt(c, "id")
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = ctrl.service.delete(id)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (ctrl Controller) createEquipmentFromOrder(c *gin.Context) {
	registrationCode, err := url.ParseToString(c, "registrationCode")
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var order models.Order
	if err := c.BindJSON(&order); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	equipment, err := ctrl.service.fulfill(order, registrationCode)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, equipment)
}
