package models

import "time"

type DbChange struct {
	ID          uint64 `gorm:"primarykey"`
	CreatedAt   time.Time
	MemberId    uint64
	EquipmentId uint64
	OrderId     uint64
	Action      string
	UserId      uint64
}

func (DbChange) TableName() string {
	return "changes"
}

type Change struct {
	ID          uint64
	CreatedAt   time.Time
	MemberId    uint64
	EquipmentId uint64
	OrderId     uint64
	Action      string
	UserId      uint64
}

const (
	OrderEquipment   string = "order-equipment"
	DeleteOrder             = "delete-order"
	UpdateOrder             = "update-order"
	OrderToEquipment        = "order-to-equipment"
	CreateMember            = "create-member"
	UpdateMember            = "update-member"
	DeleteMember            = "delete-member"
	CreateEquipment         = "create-equipment"
	DeleteEquipment         = "delete-equipment"
)

func (dbc DbChange) FromDB() *Change {
	return &Change{
		ID:          dbc.ID,
		CreatedAt:   dbc.CreatedAt,
		MemberId:    dbc.MemberId,
		EquipmentId: dbc.EquipmentId,
		OrderId:     dbc.OrderId,
		Action:      dbc.Action,
		UserId:      dbc.UserId,
	}
}

func (c Change) ToDB() DbChange {
	return DbChange{
		ID:          c.ID,
		CreatedAt:   c.CreatedAt,
		MemberId:    c.MemberId,
		EquipmentId: c.EquipmentId,
		OrderId:     c.OrderId,
		Action:      c.Action,
		UserId:      c.UserId,
	}

}
