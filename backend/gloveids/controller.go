package gloveids

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GloveIdService interface {
	GetNextGloveId() (string, error)
	MarkGloveIdAsUsed(gloveId string) error
	AddGloveId(gloveId string) error
	DeleteGloveId(gloveId string) error
}

type Controller struct {
	service GloveIdService
}

func NewController(baseRoute *gin.RouterGroup, service GloveIdService) {
	ctrl := Controller{
		service: service,
	}

	gloveIdRoute := baseRoute.Group("/glove-ids")
	{
		gloveIdRoute.GET("/next", ctrl.getNextGloveId)
		gloveIdRoute.POST("/mark-used/:id", ctrl.markGloveIdAsUsed)
		gloveIdRoute.POST("/", ctrl.addGloveId)
		gloveIdRoute.DELETE("/:id", ctrl.deleteGloveId)
	}
}

func (ctrl Controller) getNextGloveId(c *gin.Context) {
	nextId, err := ctrl.service.GetNextGloveId()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"nextId": nextId,
	})
}

func (ctrl Controller) markGloveIdAsUsed(c *gin.Context) {
	gloveId := c.Param("id")

	err := ctrl.service.MarkGloveIdAsUsed(gloveId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

type addGloveIdRequest struct {
	GloveId string `json:"gloveId" binding:"required"`
}

func (ctrl Controller) addGloveId(c *gin.Context) {
	var req addGloveIdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := ctrl.service.AddGloveId(req.GloveId); err != nil {
		c.AbortWithError(http.StatusConflict, err)
		return
	}

	c.Status(http.StatusCreated)
}

func (ctrl Controller) deleteGloveId(c *gin.Context) {
	gloveId := c.Param("id")

	if err := ctrl.service.DeleteGloveId(gloveId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}
