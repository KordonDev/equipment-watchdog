package gloveids

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GloveIdService interface {
	GetNextGloveId() (string, error)
	MarkGloveIdAsUsed(gloveId string) error
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
