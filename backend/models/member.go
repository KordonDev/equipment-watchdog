package models

import (
	"time"
)

type DbMember struct {
	ID        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"unique"`
	Group     string
	Equipment []*DbEquipment `gorm:"foreignKey:MemberID"`
}

func (DbMember) TableName() string {
	return "members"
}

type Member struct {
	Id        uint64                       `json:"id"`
	Name      string                       `json:"name"`
	Group     string                       `json:"group"`
	Equipment map[EquipmentType]*Equipment `json:"equipments"`
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
	e := make(map[EquipmentType]*Equipment, 0)
	for _, equipment := range dbm.Equipment {
		e[equipment.Type] = equipment.FromDB()
	}
	return &Member{
		Id:        dbm.ID,
		Name:      dbm.Name,
		Group:     dbm.Group,
		Equipment: e,
	}
}

func (m Member) ListToMap(equipments []*Equipment, memberId uint64) map[EquipmentType]*Equipment {
	e := make(map[EquipmentType]*Equipment, 0)
	for _, equipment := range equipments {
		equipment.MemberID = memberId
		e[equipment.Type] = equipment
	}
	return e
}
