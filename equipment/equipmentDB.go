package equipment

import (
	"gorm.io/gorm"
)

type equipmentDB struct {
	db *gorm.DB
}

func NewEquipmentDB(db *gorm.DB) *equipmentDB {
	db.AutoMigrate(&dbEquipment{})

	return &equipmentDB{
		db: db,
	}
}

func (edb *equipmentDB) getById(id uint64) (*equipment, error) {
	var e dbEquipment
	err := edb.db.Model(&dbEquipment{}).First(&e, "ID = ?", id).Error

	if err != nil {
		return &equipment{}, err
	}

	return e.fromDB(), nil
}

func (edb *equipmentDB) getByType(equipmentType string) ([]*equipment, error) {
	var dbEquipment []dbEquipment

	err := edb.db.Where("type = ?", equipmentType).Find(dbEquipment).Error
	if err != nil {
		return make([]*equipment, 0), err
	}

	equipment := make([]*equipment, 0)
	for _, v := range dbEquipment {
		equipment = append(equipment, v.fromDB())
	}

	return equipment, nil
}

func (edb *equipmentDB) Create(equipment *equipment) (*equipment, error) {
	e := equipment.toDb()
	err := edb.db.Create(&e).Error
	if err != nil {
		return nil, err
	}
	return e.fromDB(), nil
}
