package models

import "time"

type DbChange struct { // fixme give table good name
	ID        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	ToMember  uint64
	Equipment uint64
	Order     uint64 // Fixme rename Order is keyword
	Action    string
	ByUser    uint64
}

type Change struct {
	ID        uint64
	CreatedAt time.Time
	ToMember  uint64
	Equipment uint64
	Order     uint64
	Action    string
	ByUser    uint64
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
		ID:        dbc.ID,
		CreatedAt: dbc.CreatedAt,
		ToMember:  dbc.ToMember,
		Equipment: dbc.Equipment,
		Order:     dbc.Order,
		Action:    dbc.Action,
		ByUser:    dbc.ByUser,
	}
}

func (c Change) ToDB() DbChange {
	return DbChange{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		ToMember:  c.ToMember,
		Equipment: c.Equipment,
		Order:     c.Order,
		Action:    c.Action,
		ByUser:    c.ByUser,
	}

}
