package members

import (
	"fmt"

	"github.com/kordondev/equipment-watchdog/audit"
	"github.com/kordondev/equipment-watchdog/models"

	"gorm.io/gorm"
)

type memberDB struct {
	*gorm.DB
}

func NewMemberDB(db *gorm.DB) *memberDB {
	return &memberDB{
		DB: db,
	}
}

func (mdb *memberDB) getAllMember() ([]*models.Member, error) {
	var dbMembers []models.DbMember

	err := mdb.Preload("Equipment").Find(&dbMembers).Error

	if err != nil {
		return nil, err
	}

	members := make([]*models.Member, 0)
	for _, m := range dbMembers {
		members = append(members, m.FromDB())
	}
	return members, nil
}

func (mdb *memberDB) getMemberByName(name string) (*models.Member, error) {
	var m models.DbMember
	err := mdb.Preload("Equipment").Model(&models.DbMember{}).First(&m, "name = ?", name).Error

	if err != nil {
		return &models.Member{}, fmt.Errorf("error getting user: %s", name)
	}

	return m.FromDB(), nil
}

func (mdb *memberDB) getMemberById(id uint64) (*models.Member, error) {
	var m models.DbMember
	err := mdb.Preload("Equipment").Model(&models.DbMember{}).First(&m, "ID = ?", id).Error

	if err != nil {
		return &models.Member{}, fmt.Errorf("error getting user by id: %d", id)
	}

	return m.FromDB(), nil
}

func (mdb *memberDB) saveMember(member *models.Member) error {
	dbm := member.ToDB()

	// This is the single most dangerous call in the codebase: with a
	// nullable has-many FK, GORM's Association.Replace will set
	// member_id = NULL on every equipment row not present in dbm.Equipment.
	// We log the slice we are about to pass so we can correlate any sudden
	// equipment loss with the actual saveMember invocation that caused it.
	ids := make([]uint64, 0, len(dbm.Equipment))
	for _, e := range dbm.Equipment {
		if e == nil {
			continue
		}
		ids = append(ids, e.ID)
	}
	audit.Log("memberDB.saveMember", "system",
		audit.F("memberId", dbm.ID),
		audit.F("name", dbm.Name),
		audit.F("equipmentCount", len(dbm.Equipment)),
		audit.F("equipmentIds", fmt.Sprintf("%v", ids)),
	)

	err := mdb.Save(dbm).Error
	if err != nil {
		audit.Log("memberDB.saveMember.saveFailed", "system",
			audit.F("memberId", dbm.ID),
			audit.F("error", err.Error()),
		)
	}
	if assocErr := mdb.Model(dbm).Association("Equipment").Replace(dbm.Equipment); assocErr != nil {
		audit.Log("memberDB.saveMember.replaceFailed", "system",
			audit.F("memberId", dbm.ID),
			audit.F("error", assocErr.Error()),
		)
	}
	return err
}

func (mdb *memberDB) createMember(member *models.Member) (*models.Member, error) {
	m := member.ToDB()
	err := mdb.Create(&m).Error
	if err != nil {
		return nil, err
	}
	return m.FromDB(), nil
}

func (mdb *memberDB) deleteMember(member *models.Member) error {
	return mdb.Delete(member.ToDB()).Error
}

func (mdb *memberDB) getForIds(ids []uint64) ([]*models.Member, error) {
	dbMember := make([]models.DbMember, 0)

	err := mdb.Where("ID IN ?", ids).Find(&dbMember).Error
	if err != nil {
		return make([]*models.Member, 0), err
	}

	return listFromDB(dbMember), nil
}

func listFromDB(dbMember []models.DbMember) []*models.Member {
	member := make([]*models.Member, 0)
	for _, v := range dbMember {
		member = append(member, v.FromDB())
	}

	return member
}
