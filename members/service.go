package members

import (
	"github.com/cloudflare/cfssl/log"
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

func (s MemberService) updateMember(id uint64, um *models.Member) ([]uint64, error) {
	eqIds := make([]uint64, 0)

	for _, eT := range models.GroupWithEquipment[um.Group] {
		if um.Equipment[eT] != nil {
			eqIds = append(eqIds, um.Equipment[eT].Id)
		}
	}

	equipments, err := s.equipmentService.GetForIds(eqIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	changeIds := s.diffEquipment(id, um)
	um.Id = id
	um.Equipment = um.ListToMap(equipments, um.Id)
	return changeIds, s.db.saveMember(um)
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

func (s MemberService) diffEquipment(id uint64, nm *models.Member) []uint64 {
	om, _ := s.getMemberById(id)
	changeIds := make([]uint64, 0)
	if om != nil {
		for _, eT := range models.GroupWithEquipment[om.Group] {
			oldId := uint64(0)
			if om.Equipment[eT] != nil {
				oldId = om.Equipment[eT].Id
			}
			newId := uint64(0)
			if nm.Equipment[eT] != nil {
				newId = nm.Equipment[eT].Id
			}
			if oldId != newId {
				changeIds = append(changeIds, newId)
			}
		}
	}

	return changeIds
}
