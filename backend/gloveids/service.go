package gloveids

import (
	"gorm.io/gorm"
)

type Database interface {
	getNextAvailableId() (string, error)
	markIdAsUsed(gloveId string) error
}

type Service struct {
	db Database
}

func NewGloveIdService(database *gorm.DB) *Service {
	db := newGloveIdDB(database)
	return &Service{
		db: db,
	}
}

func (s *Service) GetNextGloveId() (string, error) {
	return s.db.getNextAvailableId()
}

func (s *Service) MarkGloveIdAsUsed(gloveId string) error {
	return s.db.markIdAsUsed(gloveId)
}
