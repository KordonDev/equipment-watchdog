package orders

import (
	"log"

	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

type orderDB struct {
	*gorm.DB
}

func newOrderDB(db *gorm.DB) *orderDB {
	err := db.AutoMigrate(&models.DBOrder{})
	if err != nil {
		log.Fatal(err)
	}

	return &orderDB{
		db,
	}
}

func (odb orderDB) getById(id uint64) (models.Order, error) {
	var o models.DBOrder
	err := odb.Model(&models.DBOrder{}).First(&o, "ID = ?", id).Error

	if err != nil {
		return models.Order{}, err
	}
	return o.FromDB(), nil
}

func (odb orderDB) getForMember(id uint64) ([]models.Order, error) {
  orders := make([]models.DBOrder, 0)
	err := odb.Where("member_id = ? and fulfilled_at =  \"0001-01-01 00:00:00+00:00\"", id).Find(&orders).Error

	if err != nil {
		return nil, err
	}

  result := make([]models.Order, 0)
  for _, o := range(orders) {
    result = append(result, o.FromDB())
  }


	return result, nil
}

func (odb orderDB) create(order models.Order) (models.Order, error) {
	o := order.ToDB()
	err := odb.Create(&o).Error
	if err != nil {
		return models.Order{}, err
	}
	return o.FromDB(), nil
}

func (odb orderDB) save(order *models.Order) error {
	o := order.ToDB()
	err := odb.Save(&o).Error
	return err
}

func (odb orderDB) delete(id uint64) error {
	return odb.Delete(&models.DBOrder{}, id).Error
}

func (odb orderDB) getAll(fulfilled bool) ([]models.Order, error) {
	var err error
	var result []models.DBOrder
	if fulfilled == true {
		result, err = odb.getAllFulfilled()
	} else {
		result, err = odb.getAllOpen()
	}

	if err != nil {
		return nil, err
	}

	orders := make([]models.Order, 0)
	for _, order := range result {
		orders = append(orders, order.FromDB())
	}

	return orders, nil
}

func (odb orderDB) getAllOpen() ([]models.DBOrder, error) {
	var result []models.DBOrder
	err := odb.Model(&models.DBOrder{}).Find(&result, "fulfilled_at = \"0001-01-01 00:00:00+00:00\"").Error
	return result, err
}

func (odb orderDB) getAllFulfilled() ([]models.DBOrder, error) {
	var result []models.DBOrder
	err := odb.Model(&models.DBOrder{}).Find(&result, "fulfilled_at !=  \"0001-01-01 00:00:00+00:00\"").Error
	return result, err
}
