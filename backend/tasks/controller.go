package tasks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/models"
)

type Service interface {
	getTasksByGroup(string) ([]*models.Task, error)
}

type Controller struct {
	service Service
}

func NewController(baseRoute *gin.RouterGroup, service Service) {
	ctrl := Controller{
		service: service,
	}

	tasksRoute := baseRoute.Group("/tasks")
	{
		tasksRoute.GET("/group/:group", ctrl.getTasksByGroup)
	}
}

func (ctrl Controller) getTasksByGroup(c *gin.Context) {
	group := c.Param("group")

	tasks, err := ctrl.service.getTasksByGroup(group)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, tasks)
}
