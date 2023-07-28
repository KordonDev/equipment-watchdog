package orders

import (
	"github.com/kordondev/equipment-watchdog/models"
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

func (s OrderService) create(o models.Order) (models.Order, error) {
	return s.db.create(o)
}

func (s OrderService) getById(id uint64) (models.Order, error) {
	return s.db.getById(id)
}

func (s OrderService) update(id uint64, update models.Order) (models.Order, error) {
	existing, err := s.db.getById(id)
	if err != nil {
		return models.Order{}, err
	}

	update.ID = existing.ID
	update.CreatedAt = existing.CreatedAt
	err = s.db.save(&update)
	if err != nil {
		return models.Order{}, err
	}

	return update, nil
}

func (s OrderService) delete(id uint64) error {
	return s.db.delete(id)
}

func (s OrderService) getAll(fulfilled bool) ([]models.Order, error) {
	return s.db.getAll(fulfilled)
}
