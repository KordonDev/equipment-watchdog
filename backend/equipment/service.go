package equipment

import (
	"errors"

	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

type EquipmentDatabase interface {
	CreateEquipment(*models.Equipment) (*models.Equipment, error)
	getById(uint64) (*models.Equipment, error)
	getAll() ([]*models.Equipment, error)
	getByType(string) ([]*models.Equipment, error)
	delete(uint64) error
	getForIds([]uint64) ([]*models.Equipment, error)
	getFreeEquipment() ([]*models.Equipment, error)
	save(*models.Equipment) (*models.Equipment, error)
	getByMemberIdAndType(uint64, models.EquipmentType) (*models.Equipment, error)
	registrationCodeExists(string) bool
	getByRegistrationCode(string) (*models.Equipment, error)
}

type GloveIdService interface {
	MarkGloveIdAsUsed(gloveId string) error
}

type EquipmentService struct {
	db             EquipmentDatabase
	gloveIdService GloveIdService
}

func NewEquipmentService(db *gorm.DB, gloveIdService GloveIdService) EquipmentService {
	return EquipmentService{
		db:             newEquipmentDB(db),
		gloveIdService: gloveIdService,
	}
}

func (s EquipmentService) getEquipmentById(id uint64) (*models.Equipment, error) {
	return s.db.getById(id)
}

func (s EquipmentService) getAllEquipmentByType(eType string) ([]*models.Equipment, error) {
	return s.db.getByType(eType)
}

func (s EquipmentService) createEquipment(e models.Equipment) (*models.Equipment, error) {
	equip, err := s.db.CreateEquipment(&e)
	if err != nil {
		return nil, err
	}
	if equip.Type == models.Gloves {
		err = s.gloveIdService.MarkGloveIdAsUsed(equip.RegistrationCode)
	}
	return equip, nil
}

func (s EquipmentService) deleteEquipment(id uint64) error {
	return s.db.delete(id)
}

func (s EquipmentService) GetForIds(ids []uint64) ([]*models.Equipment, error) {
	return s.db.getForIds(ids)
}

func (s EquipmentService) getAllEquipment() ([]*models.Equipment, error) {
	return s.db.getAll()
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

func (s EquipmentService) ReplaceEquipmentForMember(equipment models.Equipment) (*models.Equipment, *models.Equipment, error) {
	newEquipment, err := s.getEquipmentById(equipment.Id)
	if err != nil {
		return nil, nil, err
	}

	oldEquipment, _ := s.db.getByMemberIdAndType(equipment.MemberID, equipment.Type)
	if oldEquipment != nil {
		oldEquipment.MemberID = 0
		_, err := s.save(*oldEquipment)
		if err != nil {
			return nil, nil, err
		}
	}

	newEquipment.MemberID = equipment.MemberID
	e, err := s.save(*newEquipment)
	if err != nil {
		return nil, nil, err
	}

	return e, oldEquipment, nil
}

func (s EquipmentService) RegistrationCodeExists(rc string) bool {
	return s.db.registrationCodeExists(rc)
}

func (s EquipmentService) AssignOrCreateEquipmentForMember(memberId uint64, equipment models.Equipment) (*models.Equipment, *models.Equipment, error) {
	existingEquipment, err := s.db.getByRegistrationCode(equipment.RegistrationCode)
	var newEquipment *models.Equipment
	if errors.Is(err, gorm.ErrRecordNotFound) {
		equipment.MemberID = memberId
		newEquipment, err = s.createEquipment(equipment)
	}

	oldEquipment, _ := s.db.getByMemberIdAndType(memberId, equipment.Type)
	if oldEquipment != nil && oldEquipment.Id != equipment.Id {
		if _, err := s.UnassignEquipment(oldEquipment.Id); err != nil {
			return nil, nil, err
		}
	}

	if newEquipment != nil {
		return newEquipment, oldEquipment, nil
	}

	existingEquipment.MemberID = memberId
	existingEquipment.Size = equipment.Size
	saved, err := s.save(*existingEquipment)
	return saved, oldEquipment, err
}

func (s EquipmentService) UnassignEquipment(equipmentId uint64) (*models.Equipment, error) {
	equipment, err := s.db.getById(equipmentId)
	if err != nil {
		return nil, err
	}

	if equipment.Type == models.Helmet {
		return equipment, s.deleteEquipment(equipmentId)
	}

	equipment.MemberID = 0
	return s.save(*equipment)
}
