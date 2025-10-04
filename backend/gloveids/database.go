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
	usedGloveIds := make([]uint64, 0)
	err := gdb.Model(&models.DbGloveId{}).Where("used = true").Pluck("glove_id", &usedGloveIds).Error
	if err != nil {
		return "", err
	}
	// make list to map for faster lookup
	usedGloveIdMap := make(map[uint64]bool)
	for _, id := range usedGloveIds {
		usedGloveIdMap[id] = true
	}

	for id := 1; id < 1000; id++ {
		if usedGloveIdMap[uint64(id)] {
			return fmt.Sprintf("%d", id), nil
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

	if existingId.Used {
		return errors.New("glove ID already used")
	}

	return gdb.Model(&existingId).Update("used", true).Error
}
