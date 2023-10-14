package changes

import (
	"github.com/cloudflare/cfssl/log"

	"github.com/kordondev/equipment-watchdog/models"

	"gorm.io/gorm"
)

type changeDB struct {
	*gorm.DB
}

func newChangeDB(db *gorm.DB) *changeDB {
	err := db.AutoMigrate(&models.DbChange{})
	if err != nil {
		log.Fatal(err)
	}

	return &changeDB{
		DB: db,
	}
}

func (mdb *changeDB) getAllChanges() ([]*models.Change, error) {
	var dbChanges []models.DbChange

	err := mdb.Find(&dbChanges).Error
	if err != nil {
		return nil, err
	}

	return listFromDB(dbChanges), nil
}

func (mdb *changeDB) getForEquipment(id uint64) ([]*models.Change, error) {
	var dbChanges []models.DbChange

	err := mdb.Where("equipment == ?", id).Find(&dbChanges).Error
	if err != nil {
		return nil, err
	}

	return listFromDB(dbChanges), nil
}

func (mdb *changeDB) getForOrder(id uint64) ([]*models.Change, error) {
	var dbChanges []models.DbChange

	err := mdb.Where("order == ?", id).Find(&dbChanges).Error
	if err != nil {
		return nil, err
	}

	return listFromDB(dbChanges), nil
}

func (mdb *changeDB) getForMember(id uint64) ([]*models.Change, error) {
	var dbChanges []models.DbChange

	err := mdb.Where("to_member == ?", id).Find(&dbChanges).Error
	if err != nil {
		return nil, err
	}

	return listFromDB(dbChanges), nil
}

func (mdb *changeDB) save(change models.Change) (*models.Change, error) {
	c := change.ToDB()
	err := mdb.Save(&c).Error
	if err != nil {
		return nil, err
	}
	return c.FromDB(), nil
}

func listFromDB(dbChanges []models.DbChange) []*models.Change {
	changes := make([]*models.Change, 0)
	for _, c := range dbChanges {
		changes = append(changes, c.FromDB())
	}
	return changes
}
