package members

import (
	"fmt"

	"gorm.io/gorm"
)

type memberDB struct {
	db *gorm.DB
}

func NewMemberDB(db *gorm.DB) *memberDB {
	db.AutoMigrate(&dbMember{})

	return &memberDB{
		db: db,
	}
}

func (mdb *memberDB) GetAllMember() ([]*member, error) {
	var dbMembers []dbMember

	err := mdb.db.Find(&dbMembers).Error

	if err != nil {
		return nil, err
	}

	var members []*member
	for _, m := range dbMembers {
		members = append(members, m.fromDB())
	}
	return members, nil
}

func (mdb *memberDB) GetMemberByName(name string) (*member, error) {
	var m dbMember
	err := mdb.db.Model(&dbMember{}).First(&m, "name = ?", name).Error

	if err != nil {
		return &member{}, fmt.Errorf("error getting user: %s", name)
	}

	return m.fromDB(), nil
}

func (mdb *memberDB) GetMemberById(id uint64) (*member, error) {
	var m dbMember
	err := mdb.db.Model(&dbMember{}).First(&m, "ID = ?", id).Error

	if err != nil {
		return &member{}, fmt.Errorf("error getting user by id: %d", id)
	}

	return m.fromDB(), nil
}

func (mdb *memberDB) SaveMember(member *member) error {
	return mdb.db.Save(member.toDB()).Error
}

func (mdb *memberDB) CreateMember(member *member) (*member, error) {
	err := mdb.db.Create(member.toDB()).Error
	if err != nil {
		return nil, err
	}
	return mdb.GetMemberByName(member.Name)
}

func (mdb *memberDB) DeleteMember(member *member) error {
	return mdb.db.Delete(member.toDB()).Error
}
