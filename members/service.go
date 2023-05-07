package members

import (
	"github.com/cloudflare/cfssl/log"
	"github.com/kordondev/equipment-watchdog/equipment"
	"github.com/kordondev/equipment-watchdog/models"
)

type MemberDatabase interface {
	GetMemberById(id uint64) (*models.Member, error)
	GetAllMember() ([]*models.Member, error)
	DeleteMember(*models.Member) error
	CreateMember(*models.Member) (*models.Member, error)
	SaveMember(*models.Member) error
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

func (s MemberService) GetAllMembers() ([]*models.Member, error) {
	return s.db.GetAllMember()
}

func (s MemberService) GetMemberById(id uint64) (*models.Member, error) {
	return s.db.GetMemberById(id)
}

func (s MemberService) UpdateMember(id uint64, um *models.Member) error {
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
	return s.db.SaveMember(um)
}

func (s MemberService) CreateMember(m *models.Member) (*models.Member, error) {
	return s.db.CreateMember(m)
}

func (s MemberService) DeleteMemberById(id uint64) error {
	return s.db.DeleteMember(&models.Member{Id: id})
}

func (s MemberService) GetAllGroups() map[models.Group][]models.EquipmentType {
	return models.GroupWithEquipment
}
