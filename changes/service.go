package changes

import (
	"fmt"
	"sort"
	"time"

	"github.com/cloudflare/cfssl/log"
	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

type ChangeDatabase interface {
	getAllChanges() ([]*models.Change, error)
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

func (cs ChangeService) getAll() ([]*models.Change, error) {
	chs, err := cs.db.getAllChanges()
	if err != nil {
		return nil, err
	}

	log.Infof("changes: v", cs.enrich(chs))

	return chs, nil
}

func (cs ChangeService) enrich(chs []*models.Change) []string {
	changes := make([]string, len(chs))

	eids := make([]uint64, 0)
	for _, c := range chs {
		if c.Equipment != 0 {
			eids = append(eids, c.Equipment)
		}
	}

	uids := make([]uint64, 0)
	for _, c := range chs {
		if c.ByUser != 0 {
			uids = append(uids, c.ByUser)
		}
	}

	oids := make([]uint64, 0)
	for _, c := range chs {
		if c.Order != 0 {
			oids = append(oids, c.Order)
		}
	}

	mids := make([]uint64, 0)
	for _, c := range chs {
		if c.ToMember != 0 {
			mids = append(mids, c.ToMember)
		}
	}

	eqs, _ := cs.equipmentService.GetForIds(eids)
	uss, _ := cs.userService.GetForIds(uids)
	mes, _ := cs.memberService.GetForIds(mids)
	ors, _ := cs.orderService.GetForIds(oids)
	var msg string

	for _, c := range chs {

		e := getEquipmentMessage(eqs, c.Equipment)
		m := getMemberMessage(mes, c.ToMember)
		u := getUserMessage(uss, c.ByUser)
		o := getOrderMessage(ors, c.Order)
		t := getTimeMessage(c.CreatedAt)

		switch c.Action {
		case models.UpdateMember:
			msg = fmt.Sprintf("Ausrüstung (%v) vergeben an %v durch %v (%v)\n", e, m, u, t)
		case models.OrderEquipment:
			msg = fmt.Sprintf("Bestellung %v erstellt von %v (%v)\n", o, u, t)
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
		changes = append(changes, msg)
	}

	return changes
}

func getEquipmentMessage(eqs []*models.Equipment, eId uint64) string {
	idx := sort.Search(len(eqs), func(i int) bool {
		return eqs[i].Id == eId
	})
	if idx >= 0 && idx < len(eqs) {
		return fmt.Sprintf("%v (%v - %v)", eqs[idx].Type, eqs[idx].Size, eqs[idx].RegistrationCode)
	}
	log.Warningf("Equipment id %v nof found in %+v", eId, eqs)
	return fmt.Sprintf("id %v", eId)
}

func getUserMessage(uss []*models.User, uId uint64) string {
	idx := sort.Search(len(uss), func(i int) bool {
		return uss[i].ID == uId
	})
	if idx >= 0 && idx < len(uss) {
		return fmt.Sprintf("Nutzer %v", uss[idx].Name)
	}
	return fmt.Sprintf("Nutzer id %v", uId)
}

func getMemberMessage(mes []*models.Member, mId uint64) string {
	idx := sort.Search(len(mes), func(i int) bool {
		return mes[i].Id == mId
	})
	if idx >= 0 && idx < len(mes) {
		return fmt.Sprintf("Mitglied %v", mes[idx].Name)
	}
	return fmt.Sprintf("Mitglied id %v", mId)
}

func getOrderMessage(ors []models.Order, oId uint64) string {
	idx := sort.Search(len(ors), func(i int) bool {
		return ors[i].ID == oId
	})
	if idx >= 0 && idx < len(ors) {
		return fmt.Sprintf("Bestellung %v (%v)", ors[idx].Type, ors[idx].Size)
	}
	return fmt.Sprintf("Bestellung id %v", oId)
}

func getTimeMessage(t time.Time) string {
	//return t.Format("Mon 01.02.2006 15:45 Uhr")
	return time.Now().Format("Mon 01.02.2006 15:45")
}
