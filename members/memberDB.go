package members

import (
	"fmt"
	"log"

	"github.com/kordondev/equipment-watchdog/models"

	"gorm.io/gorm"
)

type memberDB struct {
	db *gorm.DB
}

func newMemberDB(db *gorm.DB) *memberDB {
	err := db.AutoMigrate(&models.DbMember{})
	if err != nil {
		log.Fatal(err)
	}

	return &memberDB{
		db: db,
	}
}

func (mdb *memberDB) GetAllMember() ([]*models.Member, error) {
	var dbMembers []models.DbMember

	err := mdb.db.Preload("Equipment").Find(&dbMembers).Error

	if err != nil {
		return nil, err
	}

	members := make([]*models.Member, 0)
	for _, m := range dbMembers {
		members = append(members, m.FromDB())
	}
	return members, nil
}

func (mdb *memberDB) GetMemberByName(name string) (*models.Member, error) {
	var m models.DbMember
	err := mdb.db.Preload("Equipment").Model(&models.DbMember{}).First(&m, "name = ?", name).Error

	if err != nil {
		return &models.Member{}, fmt.Errorf("error getting user: %s", name)
	}

	return m.FromDB(), nil
}

func (mdb *memberDB) GetMemberById(id uint64) (*models.Member, error) {
	var m models.DbMember
	err := mdb.db.Preload("Equipment").Model(&models.DbMember{}).First(&m, "ID = ?", id).Error

	if err != nil {
		return &models.Member{}, fmt.Errorf("error getting user by id: %d", id)
	}

	return m.FromDB(), nil
}

func (mdb *memberDB) SaveMember(member *models.Member) error {
	dbm := member.ToDB()
	err := mdb.db.Save(dbm).Error
	mdb.db.Model(dbm).Association("Equipment").Replace(dbm.Equipment)
	return err
}

func (mdb *memberDB) CreateMember(member *models.Member) (*models.Member, error) {
	m := member.ToDB()
	err := mdb.db.Create(&m).Error
	if err != nil {
		return nil, err
	}
	return m.FromDB(), nil
}

func (mdb *memberDB) DeleteMember(member *models.Member) error {
	return mdb.db.Delete(member.ToDB()).Error
}
