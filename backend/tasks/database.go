package tasks

import (
	"fmt"

	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

type taskDB struct {
	*gorm.DB
}

func newTaskDB(db *gorm.DB) *taskDB {
	return &taskDB{
		DB: db,
	}
}

func (tdb *taskDB) getTasksByGroup(group string) ([]*models.Task, error) {
	var dbTasks []models.DbTask

	err := tdb.Model(&models.DbTask{}).Find(&dbTasks, "`group` = ?", group).Error

	if err != nil {
		return nil, fmt.Errorf("error getting tasks for group: %s", group)
	}

	tasks := make([]*models.Task, 0, len(dbTasks))
	for _, t := range dbTasks {
		tasks = append(tasks, t.ToTask())
	}
	return tasks, nil
}
