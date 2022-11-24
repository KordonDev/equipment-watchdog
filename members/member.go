package members

import "time"

type dbMember struct {
	ID        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"unique"`
	Group     Group
}

func (dbMember) TableName() string {
	return "member"
}

type member struct {
	Id    uint64 `json:"id"`
	Name  string `json:"name"`
	Group Group  `json:"group"`
}

func (m *member) toDB() *dbMember {
	return &dbMember{
		ID:    m.Id,
		Name:  m.Name,
		Group: m.Group,
	}
}

func (dbm dbMember) fromDB() *member {
	return &member{
		Id:    dbm.ID,
		Name:  dbm.Name,
		Group: dbm.Group,
	}
}
