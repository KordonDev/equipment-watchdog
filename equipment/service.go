package equipment

import (
	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

type EquipmentDatabase interface {
	CreateEquipent(*models.Equipment) (*models.Equipment, error)
	getById(uint64) (*models.Equipment, error)
	getByType(string) ([]*models.Equipment, error)
	delete(uint64) error
	getAllByIds([]uint64) ([]*models.Equipment, error)
	getFreeEquipment() ([]*models.Equipment, error)
	save(*models.Equipment) (*models.Equipment, error)
	getByMemberIdAndType(uint64, models.EquipmentType) (*models.Equipment, error)
	registrationCodeExists(string) bool
}

type EquipmentService struct {
	db EquipmentDatabase
}

func NewEquipmentService(db *gorm.DB) EquipmentService {
	return EquipmentService{
		db: newEquipmentDB(db),
	}
}

func (s EquipmentService) getEquipmentById(id uint64) (*models.Equipment, error) {
	return s.db.getById(id)
}

func (s EquipmentService) getAllEquipmentByType(eType string) ([]*models.Equipment, error) {
	return s.db.getByType(eType)
}

func (s EquipmentService) createEquipment(e models.Equipment) (*models.Equipment, error) {
	return s.db.CreateEquipent(&e)
}

func (s EquipmentService) deleteEquipment(id uint64) error {
	return s.db.delete(id)
}

func (s EquipmentService) GetAllByIds(ids []uint64) ([]*models.Equipment, error) {
	return s.db.getAllByIds(ids)
}

func (s EquipmentService) getFreeEquipment() (map[models.EquipmentType][]*models.Equipment, error) {
	equipment, err := s.db.getFreeEquipment()

	equipments := make(map[models.EquipmentType][]*models.Equipment)
	for _, e := range equipment {
		equipments[e.Type] = append(equipments[e.Type], e)
	}
	return equipments, err
}

func (s EquipmentService) save(e models.Equipment) (*models.Equipment, error) {
	return s.db.save(&e)
}

func (s EquipmentService) CreateEquipmentFromOrder(order models.Order, registrationCode string) (*models.Equipment, error) {
	newEquipment := models.Equipment{
		Id:               0,
		Size:             order.Size,
		RegistrationCode: registrationCode,
		Type:             order.Type,
		MemberID:         order.MemberID,
	}
	return s.createEquipment(newEquipment)
}

func (s EquipmentService) ReplaceEquipmentForMember(equipment models.Equipment) (*models.Equipment, error) {
	newEquipment, err := s.getEquipmentById(equipment.Id)
	if err != nil {
		return nil, err
	}

	oldEquipment, _ := s.db.getByMemberIdAndType(equipment.MemberID, equipment.Type)
	if oldEquipment != nil {
		oldEquipment.MemberID = 0
		_, err := s.save(*oldEquipment)
		if err != nil {
			return nil, err
		}
	}

	newEquipment.MemberID = equipment.MemberID
	e, err := s.save(*newEquipment)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (s EquipmentService) RegistrationCodeExists(rc string) bool {
	return s.db.registrationCodeExists(rc)
}
