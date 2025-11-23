package registrationcodes

import (
	"errors"
	"time"

	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

type registrationCodesDB struct {
	*gorm.DB
}

func newDatabase(db *gorm.DB) *registrationCodesDB {
	return &registrationCodesDB{
		db,
	}
}

func (rdb registrationCodesDB) save(rc models.RegistrationCode) error {
	return rdb.Save(rc.ToDb()).Error
}

func (rdb registrationCodesDB) exists(ID string) bool {
	var rc models.DbRegistrationCode
	err := rdb.Model(&models.DbRegistrationCode{}).First(&rc, "ID = ?", ID).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func (rdb registrationCodesDB) deleteOutdated() error {
	return rdb.Where("reserved_until < ?", time.Now()).Delete(&models.DbRegistrationCode{}).Error
}
