package equipment

import (
	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

type EquipmentDatabase interface {
	getById(uint64) (*models.Equipment, error)
	getByType(string) ([]*models.Equipment, error)
	CreateEquipent(*models.Equipment) (*models.Equipment, error)
	delete(uint64) error
	getAllByIds([]uint64) ([]*models.Equipment, error)
	getFreeEquipment() ([]*models.Equipment, error)
}

type EquipmentService struct {
	db EquipmentDatabase
}

func NewEquipmentService(db *gorm.DB) EquipmentService {
	return EquipmentService{
		db: newEquipmentDB(db),
	}
}

func (s EquipmentService) GetEquipmentById(id uint64) (*models.Equipment, error) {
	return s.db.getById(id)
}

func (s EquipmentService) GetAllEquipmentByType(eType string) ([]*models.Equipment, error) {
	return s.db.getByType(eType)
}

func (s EquipmentService) CreateEquipment(e models.Equipment) (*models.Equipment, error) {
	return s.db.CreateEquipent(&e)
}

func (s EquipmentService) DeleteEquipment(id uint64) error {
	return s.db.delete(id)
}

func (s EquipmentService) GetAllByIds(ids []uint64) ([]*models.Equipment, error) {
	return s.db.getAllByIds(ids)
}

func (s EquipmentService) GetFreeEquipment() (map[models.EquipmentType][]*models.Equipment, error) {
	equipment, err := s.db.getFreeEquipment()

	equipments := make(map[models.EquipmentType][]*models.Equipment)
	for _, e := range equipment {
		equipments[e.Type] = append(equipments[e.Type], e)
	}
	return equipments, err
}
