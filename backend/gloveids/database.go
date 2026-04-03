package gloveids

import (
	"errors"
	"fmt"

	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gloveIdDB struct {
	*gorm.DB
}

func newGloveIdDB(db *gorm.DB) *gloveIdDB {
	return &gloveIdDB{
		DB: db,
	}
}

func (gdb *gloveIdDB) getNextAvailableId() (string, error) {
	existingIds := make([]uint64, 0)
	err := gdb.Model(&models.DbGloveId{}).Pluck("glove_id", &existingIds).Error
	if err != nil {
		return "", err
	}

	existingIdMap := make(map[uint64]bool)
	for _, id := range existingIds {
		existingIdMap[id] = true
	}

	for id := 1; id < 1000; id++ {
		if !existingIdMap[uint64(id)] {
			return fmt.Sprintf("%d", id), nil
		}
	}

	return "", errors.New("no available glove ID found")
}

func (gdb *gloveIdDB) markIdAsUsed(gloveId string) error {
	newEntry := &models.DbGloveId{
		GloveId: gloveId,
	}
	return gdb.Clauses(clause.OnConflict{DoNothing: true}).Create(newEntry).Error
}

func (gdb *gloveIdDB) addGloveId(gloveId string) error {
	var existing models.DbGloveId
	err := gdb.Where("glove_id = ?", gloveId).First(&existing).Error
	if err == nil {
		return fmt.Errorf("glove ID %s already exists", gloveId)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	newEntry := &models.DbGloveId{
		GloveId: gloveId,
	}
	return gdb.Create(newEntry).Error
}

func (gdb *gloveIdDB) deleteGloveId(gloveId string) error {
	return gdb.Where("glove_id = ?", gloveId).Delete(&models.DbGloveId{}).Error
}
