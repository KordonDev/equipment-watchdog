package changes

import (
	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

type ChangeDatabase interface {
	getAllChanges() ([]*models.Change, error)
}

type ChangeService struct {
	db ChangeDatabase
}

func NewChangeService(db *gorm.DB) ChangeService {
	return ChangeService{
		db: newChangeDB(db),
	}
}

func (cs ChangeService) getAll() ([]*models.Change, error) {
	return cs.db.getAllChanges()
}
