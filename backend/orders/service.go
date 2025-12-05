package orders

import (
	"time"

	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

type OrderDatabase interface {
	create(models.Order) (models.Order, error)
	getById(uint64) (models.Order, error)
	getForMember(uint64) ([]models.Order, error)
	save(*models.Order) error
	delete(uint64) error
	getAll(bool) ([]models.Order, error)
	getAllOpenOrUpdatedAfter(time.Time) ([]models.Order, error)
	getForIds([]uint64) ([]models.Order, error)
}

type EquipmentService interface {
	CreateEquipmentFromOrder(models.Order, string) (*models.Equipment, error)
	ReplaceEquipmentForMember(models.Equipment) (*models.Equipment, *models.Equipment, error)
}

type OrderService struct {
	db               OrderDatabase
	equipmentService EquipmentService
}

func NewOrderService(db *gorm.DB, equipmentService EquipmentService) *OrderService {
	return &OrderService{
		db:               newOrderDB(db),
		equipmentService: equipmentService,
	}
}

func (s OrderService) create(o models.Order) (models.Order, error) {
	return s.db.create(o)
}

func (s OrderService) getById(id uint64) (models.Order, error) {
	return s.db.getById(id)
}

func (s OrderService) getForMember(id uint64) ([]models.Order, error) {
	return s.db.getForMember(id)
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

func (s OrderService) GetForIds(ids []uint64) ([]models.Order, error) {
	return s.db.getForIds(ids)
}

func (s OrderService) fulfill(order models.Order, registrationCode string) (*models.Equipment, *models.Equipment, error) {
	equipment, err := s.equipmentService.CreateEquipmentFromOrder(order, registrationCode)
	if err != nil {
		return nil, nil, err
	}

	equipment.MemberID = order.MemberID
	_, oldEquip, err := s.equipmentService.ReplaceEquipmentForMember(*equipment)
	if err != nil {
		return nil, nil, err
	}

	order.FulfilledAt = time.Now()
	_, err = s.update(order.ID, order)
	if err != nil {
		return nil, nil, err
	}

	return equipment, oldEquip, err
}

func (s OrderService) getOpenOrRecentChanged() ([]models.Order, error) {
	oneYearAgo := time.Now().AddDate(-1, 0, 0)
	return s.db.getAllOpenOrUpdatedAfter(oneYearAgo)
}
