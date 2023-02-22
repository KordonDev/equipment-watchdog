package orders

import (
	"log"

	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

type orderDB struct {
	db *gorm.DB
}

func newOrderDB(db *gorm.DB) *orderDB {
	err := db.AutoMigrate(&models.DBOrder{})
	if err != nil {
		log.Fatal(err)
	}

	return &orderDB{
		db: db,
	}
}
func (odb orderDB) getById(id uint64) (models.Order, error) {
	var o models.DBOrder
	err := odb.db.Model(&models.DBOrder{}).First(&o, "ID = ?", id).Error

	if err != nil {
		return models.Order{}, err
	}
	return o.FromDB(), nil
}

func (odb orderDB) create(order models.Order) (models.Order, error) {
	o := order.ToDB()
	err := odb.db.Create(&o).Error
	if err != nil {
		return models.Order{}, err
	}
	return o.FromDB(), nil
}

func (odb orderDB) save(order *models.Order) error {
	o := order.ToDB()
	err := odb.db.Save(&o).Error
	return err
}

func (odb orderDB) delete(id uint64) error {
	return odb.db.Delete(&models.DBOrder{}, id).Error
}
