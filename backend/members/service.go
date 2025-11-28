package members

import (
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
	oldMember, _ := s.getMemberById(id)

	um.Id = id
	um.Equipment = oldMember.Equipment
	return s.db.saveMember(um)
}

func (s MemberService) createMember(m *models.Member) (*models.Member, error) {
	return s.db.createMember(m)
}

func (s MemberService) deleteMemberById(id uint64) error {
	return s.db.deleteMember(&models.Member{Id: id})
}

func (s MemberService) getAllGroups() map[models.Group][]models.EquipmentType {
	return models.GroupWithEquipment
}

func (s MemberService) GetForIds(ids []uint64) ([]*models.Member, error) {
	return s.db.getForIds(ids)
}

func (s MemberService) saveEquipmentForMember(memberId uint64, equipmentType models.EquipmentType, equipment models.Equipment) (*models.Equipment, *models.Equipment, error) {
	equipment.Type = equipmentType
	return s.equipmentService.AssignOrCreateEquipmentForMember(memberId, equipment)
}

func (s MemberService) removeEquipmentFromMember(memberId uint64, equipmentType models.EquipmentType) (*models.Equipment, error) {
	member, err := s.getMemberById(memberId)
	if err != nil {
		return nil, err
	}
	if member.Equipment[equipmentType] == nil {
		return nil, nil
	}
	return s.equipmentService.UnassignEquipment(member.Equipment[equipmentType].Id)
}
