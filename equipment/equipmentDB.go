package equipment

import "gorm.io/gorm"

type equipmentDB struct {
	db *gorm.DB
}

func (edb *equipmentDB) getById(id uint64) (*equipment,error) {
	var e dbEquipment
	err := edb.db.Model(&dbEquipment{}).First(&e, "ID = ?", id).Error

	if err != nil {
		return &equipment{}, err
	}

	return e.fromDB(), nil
}