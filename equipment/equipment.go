package equipment

import (
	"github.com/kordondev/equipment-watchdog/members"
	"time"
)

type dbEquipment struct {
	ID uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Type members.EquipmentType
}

func (dbEquipment) TableName() string {
	return "equipment"
}

type equipment struct {
	Id uint64 `json:"id"`
	Type members.EquipmentType `json:"type"`
}

func (e *equipment) toDb() *dbEquipment {
	return &dbEquipment{
		ID: e.Id,
		Type: e.Type,
	}
}

func (dbe dbEquipment) fromDB() *equipment {
	return &equipment{
		Id: dbe.ID,
		Type: dbe.Type,
	}
}