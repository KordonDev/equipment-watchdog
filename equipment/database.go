package equipment

import (
	"errors"
	"log"

	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

type equipmentDB struct {
	*gorm.DB
}

func newEquipmentDB(db *gorm.DB) *equipmentDB {
	err := db.AutoMigrate(&models.DbEquipment{})
	if err != nil {
		log.Fatal(err)
	}

	return &equipmentDB{
		DB: db,
	}
}

func (edb *equipmentDB) getById(id uint64) (*models.Equipment, error) {
	var e models.DbEquipment
	err := edb.Model(models.DbEquipment{}).First(&e, "ID = ?", id).Error

	if err != nil {
		return &models.Equipment{}, err
	}

	return e.FromDB(), nil
}

func (edb *equipmentDB) getByType(equipmentType string) ([]*models.Equipment, error) {
	dbEquipment := make([]models.DbEquipment, 0)

	err := edb.Where("type = ?", equipmentType).Find(&dbEquipment).Error
	if err != nil {
		return make([]*models.Equipment, 0), err
	}

	return listFormDB(dbEquipment), nil
}

func (edb *equipmentDB) CreateEquipent(equipment *models.Equipment) (*models.Equipment, error) {
	e := equipment.ToDb()
	err := edb.Create(&e).Error
	if err != nil {
		return nil, err
	}
	return e.FromDB(), nil
}

func (edb *equipmentDB) delete(id uint64) error {
	return edb.Delete(&models.DbEquipment{}, id).Error
}

func (edb *equipmentDB) getAllByIds(ids []uint64) ([]*models.Equipment, error) {
	dbEquipment := make([]models.DbEquipment, 0)

	err := edb.Where("id IN ?", ids).Find(&dbEquipment).Error
	if err != nil {
		return make([]*models.Equipment, 0), err
	}

	return listFormDB(dbEquipment), nil
}

func (edb *equipmentDB) getFreeEquipment() ([]*models.Equipment, error) {
	dbEquipment := make([]models.DbEquipment, 0)

	err := edb.Where("member_id IS null OR member_id is 0").Find(&dbEquipment).Error
	if err != nil {
		return make([]*models.Equipment, 0), err
	}

	return listFormDB(dbEquipment), nil
}

func (edb *equipmentDB) save(equipment *models.Equipment) (*models.Equipment, error) {
	e := equipment.ToDb()
	err := edb.Save(&e).Error
	if err != nil {
		return nil, err
	}
	return e.FromDB(), nil
}

func (edb *equipmentDB) getByMemberIdAndType(memberId uint64, eType models.EquipmentType) (*models.Equipment, error) {
	dbEquipment := models.DbEquipment{}
	err := edb.Where("type = ? AND member_id = ?", eType, memberId).Find(&dbEquipment).Error
	if err != nil {
		return nil, err
	}
	return dbEquipment.FromDB(), err
}

func listFormDB(dbEquipment []models.DbEquipment) []*models.Equipment {
	equipment := make([]*models.Equipment, 0)
	for _, v := range dbEquipment {
		equipment = append(equipment, v.FromDB())
	}

	return equipment
}

func (edb *equipmentDB) registrationCodeExists(rc string) bool {
	var e models.DbEquipment
	err := edb.Model(&models.DbEquipment{}).First(&e, "registration_code = ?", rc).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}
