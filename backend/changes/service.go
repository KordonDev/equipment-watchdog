package changes

import (
	"fmt"
	"slices"
	"time"

	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

type ChangeDatabase interface {
	getAllChanges() ([]*models.Change, error)
	getForEquipment(uint64) ([]*models.Change, error)
	getForOrder(uint64) ([]*models.Change, error)
	getForMember(uint64) ([]*models.Change, error)
}

type ChangeService struct {
	db               ChangeDatabase
	equipmentService equipmentService
	memberService    memberService
	userService      userService
	orderService     orderService
}

type equipmentService interface {
	GetForIds([]uint64) ([]*models.Equipment, error)
}
type memberService interface {
	GetForIds([]uint64) ([]*models.Member, error)
}
type orderService interface {
	GetForIds([]uint64) ([]models.Order, error)
}
type userService interface {
	GetForIds([]uint64) ([]*models.User, error)
}

func NewChangeService(db *gorm.DB, es equipmentService, ms memberService, us userService, os orderService) ChangeService {
	return ChangeService{
		db:               newChangeDB(db),
		equipmentService: es,
		memberService:    ms,
		userService:      us,
		orderService:     os,
	}
}

func (cs ChangeService) getAll() ([]string, error) {
	chs, err := cs.db.getAllChanges()
	if err != nil {
		return nil, err
	}

	return cs.enrich(chs), nil
}

func (cs ChangeService) getForEquipment(id uint64) ([]string, error) {
	chs, err := cs.db.getForEquipment(id)
	if err != nil {
		return nil, err
	}

	return cs.enrich(chs), nil
}

func (cs ChangeService) getForOrder(id uint64) ([]string, error) {
	chs, err := cs.db.getForOrder(id)
	if err != nil {
		return nil, err
	}

	return cs.enrich(chs), nil
}

func (cs ChangeService) getForMember(id uint64) ([]string, error) {
	chs, err := cs.db.getForMember(id)
	if err != nil {
		return nil, err
	}

	return cs.enrich(chs), nil
}

func (cs ChangeService) enrich(chs []*models.Change) []string {
	changes := make([]string, len(chs))

	eids := make([]uint64, 0)
	for _, c := range chs {
		if c.EquipmentId != 0 && !slices.Contains(eids, c.EquipmentId) {
			eids = append(eids, c.EquipmentId)
		}
	}

	uids := make([]uint64, 0)
	for _, c := range chs {
		if c.UserId != 0 && !slices.Contains(uids, c.UserId) {
			uids = append(uids, c.UserId)
		}
	}

	oids := make([]uint64, 0)
	for _, c := range chs {
		if c.OrderId != 0 && !slices.Contains(oids, c.OrderId) {
			oids = append(oids, c.OrderId)
		}
	}

	mids := make([]uint64, 0)
	for _, c := range chs {
		if c.MemberId != 0 && !slices.Contains(mids, c.MemberId) {
			mids = append(mids, c.MemberId)
		}
	}

	eqs, _ := cs.equipmentService.GetForIds(eids)
	uss, _ := cs.userService.GetForIds(uids)
	mes, _ := cs.memberService.GetForIds(mids)
	ors, _ := cs.orderService.GetForIds(oids)
	var msg string

	for i, c := range chs {

		e := getEquipmentMessage(eqs, c.EquipmentId)
		m := getMemberMessage(mes, c.MemberId)
		u := getUserMessage(uss, c.UserId)
		o := getOrderMessage(ors, c.OrderId)
		t := getTimeMessage(c.CreatedAt)

		switch c.Action {
		case models.UpdateMember:
			msg = fmt.Sprintf("Ausrüstung (%v) vergeben an %v durch %v (%v)\n", e, m, u, t)
		case models.CreateOrder:
			msg = fmt.Sprintf("Bestellung %v erstellt für %v von %v (%v)\n", o, m, u, t)
		case models.DeleteOrder:
			msg = fmt.Sprintf("Bestellung %v gelöscht von %v (%v)\n", o, u, t)
		case models.OrderToEquipment:
			msg = fmt.Sprintf("Bestellung %v zu %v gemacht und %v zugewiesen von %v (%v)\n", o, e, m, u, t)
		case models.CreateMember:
			msg = fmt.Sprintf("Mitglied %v erstellt von %v (%v)\n", m, u, t)
		case models.DeleteMember:
			msg = fmt.Sprintf("Mitglied %v gelöscht von %v (%v)\n", m, u, t)
		case models.CreateEquipment:
			msg = fmt.Sprintf("Mitglied %v gelöscht von %v (%v)\n", m, u, t)
		default:
			msg = c.Action
		}

		changes[i] = msg
	}

	return changes
}

func getEquipmentMessage(eqs []*models.Equipment, eId uint64) string {
	for _, eq := range eqs {
		if eq.Id == eId {
			return fmt.Sprintf("%v (%v - %v)", eq.Type, eq.Size, eq.RegistrationCode)
		}
	}
	return fmt.Sprintf("id %v", eId)
}

func getUserMessage(uss []*models.User, uId uint64) string {
	for _, us := range uss {
		if us.ID == uId {
			return fmt.Sprintf("Nutzer %v", us.Name)
		}
	}
	return fmt.Sprintf("Nutzer id %v", uId)
}

func getMemberMessage(mes []*models.Member, mId uint64) string {
	for _, me := range mes {
		if me.Id == mId {
			return fmt.Sprintf("%v", me.Name)
		}
	}
	return fmt.Sprintf("id %v", mId)
}

func getOrderMessage(ors []models.Order, oId uint64) string {
	for _, or := range ors {
		if or.ID == oId {
			return fmt.Sprintf("%v (%v)", or.Type, or.Size)
		}
	}
	return fmt.Sprintf("id %v", oId)
}

func getTimeMessage(t time.Time) string {
	return t.Format("Mon 02.01.2006 15:04")
}
