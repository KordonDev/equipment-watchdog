package models

import (
	"time"
)

type DbMember struct {
	ID        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"unique"`
	Group     Group
	Equipment []*DbEquipment `gorm:"foreignKey:MemberID"`
}

func (DbMember) TableName() string {
	return "members"
}

type Member struct {
	Id        uint64       `json:"id"`
	Name      string       `json:"name"`
	Group     Group        `json:"group"`
	Equipment []*Equipment `json:"equipment"`
}

func (m *Member) ToDB() *DbMember {
	e := make([]*DbEquipment, 0)
	for _, equipment := range m.Equipment {
		e = append(e, equipment.ToDb())
	}
	return &DbMember{
		ID:        m.Id,
		Name:      m.Name,
		Group:     m.Group,
		Equipment: e,
	}
}

func (dbm DbMember) FromDB() *Member {
	e := make([]*Equipment, 0)
	for _, equipment := range dbm.Equipment {
		e = append(e, equipment.FromDB())
	}
	return &Member{
		Id:        dbm.ID,
		Name:      dbm.Name,
		Group:     dbm.Group,
		Equipment: e,
	}
}
