package members

import (
	"github.com/cloudflare/cfssl/log"
	"github.com/kordondev/equipment-watchdog/equipment"
	"github.com/kordondev/equipment-watchdog/models"
)

type MemberDatabase interface {
	getMemberById(id uint64) (*models.Member, error)
	getAllMember() ([]*models.Member, error)
	deleteMember(*models.Member) error
	createMember(*models.Member) (*models.Member, error)
	saveMember(*models.Member) error
}
type MemberService struct {
	db               MemberDatabase
	equipmentService *equipment.EquipmentService
}

func NewMemberService(database MemberDatabase, equipmentService *equipment.EquipmentService) MemberService {
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
	eqIds := make([]uint64, 0)

	for _, eT := range models.GroupWithEquipment[um.Group] {
		if um.Equipment[eT] != nil {
			eqIds = append(eqIds, um.Equipment[eT].Id)
		}
	}
	equipments, err := s.equipmentService.GetAllByIds(eqIds)
	if err != nil {
		log.Error(err)
		return err
	}

	um.Id = id
	um.Equipment = um.ListToMap(equipments, um.Id)
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
