package equipment

import (
	"log"

	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

type equipmentDB struct {
	db *gorm.DB
}

func newEquipmentDB(db *gorm.DB) *equipmentDB {
	err := db.AutoMigrate(&models.DbEquipment{})
	if err != nil {
		log.Fatal(err)
	}

	return &equipmentDB{
		db: db,
	}
}

func (edb *equipmentDB) getById(id uint64) (*models.Equipment, error) {
	var e models.DbEquipment
	err := edb.db.Model(&models.DbEquipment{}).First(&e, "ID = ?", id).Error

	if err != nil {
		return &models.Equipment{}, err
	}

	return e.FromDB(), nil
}

func (edb *equipmentDB) getByType(equipmentType string) ([]*models.Equipment, error) {
	dbEquipment := make([]models.DbEquipment, 0)

	err := edb.db.Where("type = ?", equipmentType).Find(&dbEquipment).Error
	if err != nil {
		return make([]*models.Equipment, 0), err
	}

	return listFormDB(dbEquipment), nil
}

func (edb *equipmentDB) Create(equipment *models.Equipment) (*models.Equipment, error) {
	e := equipment.ToDb()
	err := edb.db.Create(&e).Error
	if err != nil {
		return nil, err
	}
	return e.FromDB(), nil
}

func (edb *equipmentDB) delete(id uint64) error {
	return edb.db.Delete(&models.DbEquipment{}, id).Error
}

func (edb *equipmentDB) getAllByIds(ids []uint64) ([]*models.Equipment, error) {
	dbEquipment := make([]models.DbEquipment, 0)

	err := edb.db.Where("id IN ?", ids).Find(&dbEquipment).Error
	if err != nil {
		return make([]*models.Equipment, 0), err
	}

	return listFormDB(dbEquipment), nil
}

func listFormDB(dbEquipment []models.DbEquipment) []*models.Equipment {
	equipment := make([]*models.Equipment, 0)
	for _, v := range dbEquipment {
		equipment = append(equipment, v.FromDB())
	}

	return equipment
}
