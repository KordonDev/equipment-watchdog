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
	CreateOrder      string = "create-order"
	DeleteOrder      string = "delete-order"
	UpdateOrder      string = "update-order"
	OrderToEquipment string = "order-to-equipment"
	CreateMember     string = "create-member"
	UpdateMember     string = "update-member"
	DeleteMember     string = "delete-member"
	CreateEquipment  string = "create-equipment"
	DeleteEquipment  string = "delete-equipment"
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
	return DbChange(c)
}
