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
	ID             uint64    `json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	MemberId       uint64    `json:"memberId"`
	EquipmentId    uint64    `json:"equipmentId"`
	OrderId        uint64    `json:"orderId"`
	Action         string    `json:"action"`
	UserId         uint64    `json:"userId"`
	OldEquipmentId uint64    `json:"oldEquipmentId,omitempty"`
}

const (
	CreateOrder             string = "create-order"
	DeleteOrder                    = "delete-order"
	UpdateOrder                    = "update-order"
	OrderToEquipment               = "order-to-equipment"
	CreateMember                   = "create-member"
	UpdateMember                   = "update-member"
	DeleteMember                   = "delete-member"
	CreateEquipment                = "create-equipment"
	DeleteEquipment                = "delete-equipment"
	UpdateEquipmentOnMember        = "update-equipment-on-member"
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
