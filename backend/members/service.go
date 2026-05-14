package members

import (
	"github.com/kordondev/equipment-watchdog/audit"
	"github.com/kordondev/equipment-watchdog/models"
)

type MemberDatabase interface {
	getMemberById(uint64) (*models.Member, error)
	getAllMember() ([]*models.Member, error)
	deleteMember(*models.Member) error
	createMember(*models.Member) (*models.Member, error)
	saveMember(*models.Member) error
	getForIds([]uint64) ([]*models.Member, error)
}

type EquipmentService interface {
	GetForIds([]uint64) ([]*models.Equipment, error)
	AssignOrCreateEquipmentForMember(uint64, models.Equipment) (*models.Equipment, *models.Equipment, error)
	UnassignEquipment(uint64) (*models.Equipment, error)
}
type MemberService struct {
	db               MemberDatabase
	equipmentService EquipmentService
}

func NewMemberService(database MemberDatabase, equipmentService EquipmentService) MemberService {
	return MemberService{
		db:               database,
		equipmentService: equipmentService,
	}
}

func (s MemberService) getAllMembers() ([]*models.Member, error) {
	return s.db.getAllMember()
}

func (s MemberService) getMemberById(id uint64) (*models.Member, error) {
	return s.db.getMemberById(id)
}

func (s MemberService) updateMember(id uint64, um *models.Member) error {
	oldMember, loadErr := s.getMemberById(id)
	if loadErr != nil {
		// IMPORTANT: previously this error was swallowed, which meant a
		// failed load would feed an empty equipment slice into saveMember,
		// and GORM's Association.Replace would then NULL out every piece of
		// equipment for that member. We log this loudly so we have evidence
		// next time it happens.
		audit.Log("member.update.loadOldFailed", "system",
			audit.F("memberId", id),
			audit.F("error", loadErr.Error()),
		)
	}

	beforeSnapshot := audit.EquipmentSnapshot(oldMember)
	incomingSnapshot := audit.EquipmentSnapshot(um)

	um.Id = id
	um.Equipment = oldMember.Equipment

	audit.Log("member.update", "system",
		audit.F("memberId", id),
		audit.F("name", um.Name),
		audit.F("group", um.Group),
		audit.F("equipmentBefore", beforeSnapshot),
		audit.F("equipmentInRequest", incomingSnapshot),
		audit.F("equipmentSentToSave", audit.EquipmentSnapshot(um)),
	)

	if err := s.db.saveMember(um); err != nil {
		audit.Log("member.update.saveFailed", "system",
			audit.F("memberId", id),
			audit.F("error", err.Error()),
		)
		return err
	}

	// Reload to confirm the equipment is still attached after the save.
	if after, err := s.getMemberById(id); err == nil {
		audit.Log("member.update.after", "system",
			audit.F("memberId", id),
			audit.F("equipmentAfter", audit.EquipmentSnapshot(after)),
		)
	}

	return nil
}

func (s MemberService) createMember(m *models.Member) (*models.Member, error) {
	return s.db.createMember(m)
}

func (s MemberService) deleteMemberById(id uint64) error {
	// Snapshot equipment first so we can see in the log which pieces will
	// have their member_id set to NULL by the ON DELETE SET NULL FK rule.
	if before, err := s.getMemberById(id); err == nil {
		audit.Log("member.delete", "system",
			audit.F("memberId", id),
			audit.F("name", before.Name),
			audit.F("equipmentBefore", audit.EquipmentSnapshot(before)),
		)
	} else {
		audit.Log("member.delete.unknown", "system",
			audit.F("memberId", id),
			audit.F("loadError", err.Error()),
		)
	}
	return s.db.deleteMember(&models.Member{Id: id})
}

func (s MemberService) GetForIds(ids []uint64) ([]*models.Member, error) {
	return s.db.getForIds(ids)
}

func (s MemberService) saveEquipmentForMember(memberId uint64, equipmentType models.EquipmentType, equipment models.Equipment) (*models.Equipment, *models.Equipment, error) {
	equipment.Type = equipmentType
	equipment.MemberID = memberId

	if before, err := s.getMemberById(memberId); err == nil {
		audit.Log("member.equipment.assign.before", "system",
			audit.F("memberId", memberId),
			audit.F("type", equipmentType),
			audit.F("incoming", audit.EquipmentBrief(&equipment)),
			audit.F("equipmentBefore", audit.EquipmentSnapshot(before)),
		)
	}

	saved, old, err := s.equipmentService.AssignOrCreateEquipmentForMember(memberId, equipment)

	if after, lerr := s.getMemberById(memberId); lerr == nil {
		audit.Log("member.equipment.assign.after", "system",
			audit.F("memberId", memberId),
			audit.F("type", equipmentType),
			audit.F("saved", audit.EquipmentBrief(saved)),
			audit.F("old", audit.EquipmentBrief(old)),
			audit.F("equipmentAfter", audit.EquipmentSnapshot(after)),
			audit.F("error", errString(err)),
		)
	}

	return saved, old, err
}

func (s MemberService) removeEquipmentFromMember(memberId uint64, equipmentType models.EquipmentType) (*models.Equipment, error) {
	member, err := s.getMemberById(memberId)
	if err != nil {
		audit.Log("member.equipment.remove.failed", "system",
			audit.F("stage", "loadMember"),
			audit.F("memberId", memberId),
			audit.F("type", equipmentType),
			audit.F("error", err.Error()),
		)
		return nil, err
	}
	if member.Equipment[equipmentType] == nil {
		audit.Log("member.equipment.remove.noop", "system",
			audit.F("memberId", memberId),
			audit.F("type", equipmentType),
			audit.F("equipmentBefore", audit.EquipmentSnapshot(member)),
		)
		return nil, nil
	}

	target := member.Equipment[equipmentType]
	audit.Log("member.equipment.remove", "system",
		audit.F("memberId", memberId),
		audit.F("type", equipmentType),
		audit.F("target", audit.EquipmentBrief(target)),
		audit.F("equipmentBefore", audit.EquipmentSnapshot(member)),
	)

	result, err := s.equipmentService.UnassignEquipment(target.Id)

	if after, lerr := s.getMemberById(memberId); lerr == nil {
		audit.Log("member.equipment.remove.after", "system",
			audit.F("memberId", memberId),
			audit.F("type", equipmentType),
			audit.F("equipmentAfter", audit.EquipmentSnapshot(after)),
			audit.F("error", errString(err)),
		)
	}
	return result, err
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
