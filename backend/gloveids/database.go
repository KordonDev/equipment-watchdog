package gloveids

import (
	"errors"
	"fmt"
	"log"

	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

type gloveIdDB struct {
	*gorm.DB
}

func newGloveIdDB(db *gorm.DB) *gloveIdDB {
	err := db.AutoMigrate(&models.DbGloveId{})
	if err != nil {
		log.Fatal(err)
	}

	return &gloveIdDB{
		DB: db,
	}
}

func (gdb *gloveIdDB) getNextAvailableId() (string, error) {
	for id := 1; id < 1000; id++ {
		gloveId := fmt.Sprintf("%d", id)
		var count int64
		err := gdb.Model(&models.DbGloveId{}).Where("glove_id = ? and used = true", gloveId).Count(&count).Error
		if err != nil {
			return "", err
		}

		if count == 0 {
			return gloveId, nil
		}
	}

	return "", errors.New("no available glove ID found")
}

func (gdb *gloveIdDB) markIdAsUsed(gloveId string) error {
	// Prüfe ob ID bereits existiert
	var existingId models.DbGloveId
	err := gdb.Where("glove_id = ?", gloveId).First(&existingId).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// ID existiert noch nicht, erstelle sie
			newGloveId := &models.DbGloveId{
				GloveId: gloveId,
				Used:    true,
			}
			return gdb.Create(newGloveId).Error
		}
		return err
	}

	return gdb.Model(&existingId).Update("used", true).Error
}

func (gdb *gloveIdDB) getAllUsedIds() ([]string, error) {
	var gloveIds []models.DbGloveId
	err := gdb.Where("used = ?", true).Find(&gloveIds).Error
	if err != nil {
		return nil, err
	}

	usedIds := make([]string, len(gloveIds))
	for i, gId := range gloveIds {
		usedIds[i] = gId.GloveId
	}

	return usedIds, nil
}
