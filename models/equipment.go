package models

import (
	"time"
)

type EquipmentType string

const (
	Helmet   EquipmentType = "helmet"
	Jacket   EquipmentType = "jacket"
	Gloves   EquipmentType = "gloves"
	Trousers EquipmentType = "trousers"
	Boots    EquipmentType = "boots"
	TShirt   EquipmentType = "tshirt"
)

// FIXME: move this into the equiqment package
type DbEquipment struct {
	ID               uint64 `gorm:"primarykey"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Type             EquipmentType
	RegistrationCode string `gorm:"unique"`
	MemberID         uint64
	Size             string
}

func (DbEquipment) TableName() string {
	return "equipments"
}

type Equipment struct {
	Id               uint64        `json:"id"`
	Type             EquipmentType `json:"type"`
	RegistrationCode string        `json:"registrationCode"`
	MemberID         uint64        `json:"memberId"`
	Size             string        `json:"size"`
}

func (e *Equipment) ToDb() *DbEquipment {
	return &DbEquipment{
		ID:               e.Id,
		Type:             e.Type,
		RegistrationCode: e.RegistrationCode,
		MemberID:         e.MemberID,
		Size:             e.Size,
	}
}

func (dbe *DbEquipment) FromDB() *Equipment {
	return &Equipment{
		Id:               dbe.ID,
		Type:             dbe.Type,
		RegistrationCode: dbe.RegistrationCode,
		MemberID:         dbe.MemberID,
		Size:             dbe.Size,
	}
}
