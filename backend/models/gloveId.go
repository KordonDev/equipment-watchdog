package models

import (
	"time"
)

type DbGloveId struct {
	ID        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	GloveId   string `gorm:"unique;not null"`
	Used      bool   `gorm:"default:true"`
}

func (DbGloveId) TableName() string {
	return "glove_ids"
}

type GloveId struct {
	Id      uint64 `json:"id"`
	GloveId string `json:"gloveId"`
	Used    bool   `json:"used"`
}

func (g *GloveId) ToDb() *DbGloveId {
	return &DbGloveId{
		ID:      g.Id,
		GloveId: g.GloveId,
		Used:    g.Used,
	}
}

func (dbg *DbGloveId) FromDB() *GloveId {
	return &GloveId{
		Id:      dbg.ID,
		GloveId: dbg.GloveId,
		Used:    dbg.Used,
	}
}
