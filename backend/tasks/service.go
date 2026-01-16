package tasks

import (
	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

type TaskDatabase interface {
	getTasksByGroup(string) ([]*models.Task, error)
}

type TaskService struct {
	db TaskDatabase
}

func NewTaskService(db *gorm.DB) TaskService {
	database := newTaskDB(db)
	return TaskService{
		db: database,
	}
}

func (s TaskService) getTasksByGroup(group string) ([]*models.Task, error) {
	return s.db.getTasksByGroup(group)
}
