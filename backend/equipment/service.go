package equipment

import (
	"errors"

	"github.com/kordondev/equipment-watchdog/audit"
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
		audit.Log("equipment.create.failed", "system",
			audit.F("error", err.Error()),
			audit.F("type", e.Type),
			audit.F("memberId", e.MemberID),
			audit.F("registrationCode", e.RegistrationCode),
		)
		return nil, err
	}
	if equip.Type == models.Gloves {
		err = s.gloveIdService.MarkGloveIdAsUsed(equip.RegistrationCode)
	}
	audit.Log("equipment.create", "system",
		audit.F("equipmentId", equip.Id),
		audit.F("type", equip.Type),
		audit.F("memberId", equip.MemberID),
		audit.F("size", equip.Size),
		audit.F("registrationCode", equip.RegistrationCode),
	)
	return equip, nil
}

func (s EquipmentService) deleteEquipment(id uint64) error {
	// Snapshot the equipment so we know what's being lost.
	if before, err := s.db.getById(id); err == nil {
		audit.Log("equipment.delete", "system",
			audit.F("equipmentId", id),
			audit.F("before", audit.EquipmentBrief(before)),
		)
	} else {
		audit.Log("equipment.delete.unknown", "system",
			audit.F("equipmentId", id),
			audit.F("loadError", err.Error()),
		)
	}
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
	audit.Log("equipment.fromOrder", "system",
		audit.F("orderId", order.ID),
		audit.F("memberId", order.MemberID),
		audit.F("type", order.Type),
		audit.F("size", order.Size),
		audit.F("registrationCode", registrationCode),
	)
	return s.createEquipment(newEquipment)
}

func (s EquipmentService) ReplaceEquipmentForMember(equipment models.Equipment) (*models.Equipment, *models.Equipment, error) {
	newEquipment, err := s.getEquipmentById(equipment.Id)
	if err != nil {
		audit.Log("equipment.replace.failed", "system",
			audit.F("stage", "loadNew"),
			audit.F("equipmentId", equipment.Id),
			audit.F("error", err.Error()),
		)
		return nil, nil, err
	}

	oldEquipment, _ := s.db.getByMemberIdAndType(equipment.MemberID, equipment.Type)
	if oldEquipment != nil {
		audit.Log("equipment.replace.unassignOld", "system",
			audit.F("memberId", equipment.MemberID),
			audit.F("type", equipment.Type),
			audit.F("oldEquipment", audit.EquipmentBrief(oldEquipment)),
		)
		oldEquipment.MemberID = 0
		_, err := s.save(*oldEquipment)
		if err != nil {
			audit.Log("equipment.replace.failed", "system",
				audit.F("stage", "saveOld"),
				audit.F("oldEquipmentId", oldEquipment.Id),
				audit.F("error", err.Error()),
			)
			return nil, nil, err
		}
	}

	newEquipment.MemberID = equipment.MemberID
	e, err := s.save(*newEquipment)
	if err != nil {
		audit.Log("equipment.replace.failed", "system",
			audit.F("stage", "saveNew"),
			audit.F("equipmentId", newEquipment.Id),
			audit.F("error", err.Error()),
		)
		return nil, nil, err
	}

	audit.Log("equipment.replace.done", "system",
		audit.F("memberId", equipment.MemberID),
		audit.F("type", equipment.Type),
		audit.F("newEquipment", audit.EquipmentBrief(e)),
		audit.F("oldEquipment", audit.EquipmentBrief(oldEquipment)),
	)

	return e, oldEquipment, nil
}

func (s EquipmentService) RegistrationCodeExists(rc string) bool {
	return s.db.registrationCodeExists(rc)
}

func (s EquipmentService) AssignOrCreateEquipmentForMember(memberId uint64, equipment models.Equipment) (*models.Equipment, *models.Equipment, error) {
	existingEquipment, err := s.db.getByRegistrationCode(equipment.RegistrationCode)
	oldEquipment, _ := s.db.getByMemberIdAndType(memberId, equipment.Type)

	audit.Log("equipment.assignOrCreate.start", "system",
		audit.F("memberId", memberId),
		audit.F("type", equipment.Type),
		audit.F("registrationCode", equipment.RegistrationCode),
		audit.F("existing", audit.EquipmentBrief(existingEquipment)),
		audit.F("old", audit.EquipmentBrief(oldEquipment)),
	)

	var newEquipment *models.Equipment
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newEquipment, err = s.createEquipment(equipment)
	}

	if oldEquipment != nil && oldEquipment.Id != equipment.Id {
		audit.Log("equipment.assignOrCreate.replacingOld", "system",
			audit.F("memberId", memberId),
			audit.F("type", equipment.Type),
			audit.F("oldEquipmentId", oldEquipment.Id),
			audit.F("newEquipmentId", equipment.Id),
		)
		if _, err := s.UnassignEquipment(oldEquipment.Id); err != nil {
			audit.Log("equipment.assignOrCreate.failed", "system",
				audit.F("stage", "unassignOld"),
				audit.F("oldEquipmentId", oldEquipment.Id),
				audit.F("error", err.Error()),
			)
			return nil, nil, err
		}
	}

	if newEquipment != nil {
		audit.Log("equipment.assignOrCreate.doneCreated", "system",
			audit.F("memberId", memberId),
			audit.F("newEquipment", audit.EquipmentBrief(newEquipment)),
			audit.F("oldEquipment", audit.EquipmentBrief(oldEquipment)),
		)
		return newEquipment, oldEquipment, nil
	}

	existingEquipment.MemberID = memberId
	existingEquipment.Size = equipment.Size
	saved, err := s.save(*existingEquipment)
	if err != nil {
		audit.Log("equipment.assignOrCreate.failed", "system",
			audit.F("stage", "saveExisting"),
			audit.F("equipmentId", existingEquipment.Id),
			audit.F("error", err.Error()),
		)
		return saved, oldEquipment, err
	}
	audit.Log("equipment.assignOrCreate.doneAssigned", "system",
		audit.F("memberId", memberId),
		audit.F("newEquipment", audit.EquipmentBrief(saved)),
		audit.F("oldEquipment", audit.EquipmentBrief(oldEquipment)),
	)
	return saved, oldEquipment, err
}

func (s EquipmentService) UnassignEquipment(equipmentId uint64) (*models.Equipment, error) {
	equipment, err := s.db.getById(equipmentId)
	if err != nil {
		audit.Log("equipment.unassign.failed", "system",
			audit.F("stage", "load"),
			audit.F("equipmentId", equipmentId),
			audit.F("error", err.Error()),
		)
		return nil, err
	}

	if equipment.Type == models.Helmet {
		audit.Log("equipment.unassign.deleteHelmet", "system",
			audit.F("equipmentId", equipmentId),
			audit.F("memberIdBefore", equipment.MemberID),
			audit.F("before", audit.EquipmentBrief(equipment)),
		)
		return equipment, s.deleteEquipment(equipmentId)
	}

	memberIdBefore := equipment.MemberID
	equipment.MemberID = 0
	saved, err := s.save(*equipment)
	if err != nil {
		audit.Log("equipment.unassign.failed", "system",
			audit.F("stage", "save"),
			audit.F("equipmentId", equipmentId),
			audit.F("error", err.Error()),
		)
		return saved, err
	}
	audit.Log("equipment.unassign", "system",
		audit.F("equipmentId", equipmentId),
		audit.F("type", equipment.Type),
		audit.F("memberIdBefore", memberIdBefore),
		audit.F("memberIdAfter", uint64(0)),
	)
	return saved, nil
}
