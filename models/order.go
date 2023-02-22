package models

import "time"

type DBOrder struct {
	ID        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Type      EquipmentType
	MemberID  uint64
	Size      string
}

func (DBOrder) TableName() string {
	return "orders"
}

type Order struct {
	ID        uint64        `json:"id"`
	CreatedAt time.Time     `json:"createdAt"`
	Type      EquipmentType `json:"type"`
	MemberID  uint64        `json:"memberId"`
	Size      string        `json:"size"`
}

func (o DBOrder) FromDB() Order {
	return Order{
		ID:        o.ID,
		CreatedAt: o.CreatedAt,
		Type:      o.Type,
		MemberID:  o.MemberID,
		Size:      o.Size,
	}
}

func (o Order) ToDB() DBOrder {
	return DBOrder{
		ID:        o.ID,
		Type:      o.Type,
		MemberID:  o.MemberID,
		Size:      o.Size,
		CreatedAt: o.CreatedAt,
	}
}
