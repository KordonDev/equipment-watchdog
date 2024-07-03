package members

import (
	"fmt"
	"log"

	"github.com/kordondev/equipment-watchdog/models"

	"gorm.io/gorm"
)

type memberDB struct {
	*gorm.DB
}

func NewMemberDB(db *gorm.DB) *memberDB {
	err := db.AutoMigrate(&models.DbMember{})
	if err != nil {
		log.Fatal(err)
	}

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
	err := mdb.Save(dbm).Error
	mdb.Model(dbm).Association("Equipment").Replace(dbm.Equipment)
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
