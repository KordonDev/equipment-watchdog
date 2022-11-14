package members

import (
	"fmt"

	"gorm.io/gorm"
)

type memberDB struct {
	db *gorm.DB
}

func NewMemberDB(db *gorm.DB) *memberDB {
	db.AutoMigrate(&Member{})

	return &memberDB{
		db: db,
	}
}

func (mdb *memberDB) GetMember(name string) (*Member, error) {
	var m Member
	err := mdb.db.Model(&Member{}).First(&m, "name = ?", name).Error

	if err != nil {
		return &Member{}, fmt.Errorf("error getting user: %s", name)
	}

	return &m, nil
}

func (mdb *memberDB) SaveMember(member *Member) error {
	return mdb.db.Save(member).Error
}

func (mdb *memberDB) DeleteMember(member *Member) error {
	return mdb.db.Delete(member).Error
}
