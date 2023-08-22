package registrationcodes

import (
	"github.com/cloudflare/cfssl/log"
	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

type registrationCodesDB struct {
	*gorm.DB
}

func newDatabase(db *gorm.DB) *registrationCodesDB {
	err := db.AutoMigrate(&models.DbRegistrationCode{})
	if err != nil {
		log.Fatal(err)
	}

	return &registrationCodesDB{
		db,
	}
}

func (rdb registrationCodesDB) save(rc models.RegistrationCode) error {
  return rdb.Save(rc).Error
}

func (rdb registrationCodesDB) exists(ID string) bool {
  var rc models.RegistrationCode
	err := rdb.Model(&models.DbRegistrationCode{}).First(&rc, "ID = ", ID).Error
  if err != nil {
    return true
  }
  return rc.ID == ID
}
