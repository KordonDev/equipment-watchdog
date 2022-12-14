package models

import (
	"time"
)

type DbEquipment struct {
	ID               uint64 `gorm:"primarykey"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Type             EquipmentType
	RegistrationCode string `gorm:"unique"`
	MemberID         uint64
}

func (DbEquipment) TableName() string {
	return "equipments"
}

type Equipment struct {
	Id               uint64        `json:"id"`
	Type             EquipmentType `json:"type"`
	RegistrationCode string        `json:"registrationCode"`
	MemberID         uint64        `json:"memberId"`
}

func (e *Equipment) ToDb() *DbEquipment {
	return &DbEquipment{
		ID:               e.Id,
		Type:             e.Type,
		RegistrationCode: e.RegistrationCode,
		MemberID:         e.MemberID,
	}
}

func (dbe *DbEquipment) FromDB() *Equipment {
	return &Equipment{
		Id:               dbe.ID,
		Type:             dbe.Type,
		RegistrationCode: dbe.RegistrationCode,
		MemberID:         dbe.MemberID,
	}
}
