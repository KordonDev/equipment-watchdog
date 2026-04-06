package models

import "time"

type Task struct {
	ID            uint64 `json:"id"`
	MemberID      uint64 `json:"memberId"`
	Group         string `json:"group"`
	EquipmentType string `json:"equipmentType"`
	Type          string `json:"type"`
}

type DbTask struct {
	ID            uint64 `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	MemberID      uint64 `gorm:"column:member_id"`
	EquipmentType string `gorm:"column:equipment_type"`
	Group         string `gorm:"column:group"`
	Type          string `gorm:"column:type"`
}

func (DbTask) TableName() string {
	return "tasks"
}

func (t *Task) ToDBTask() *DbTask {
	return &DbTask{
		ID:            t.ID,
		MemberID:      t.MemberID,
		Group:         t.Group,
		EquipmentType: t.EquipmentType,
		Type:          t.Type,
	}
}

func (dbt *DbTask) ToTask() *Task {
	return &Task{
		ID:            dbt.ID,
		MemberID:      dbt.MemberID,
		Group:         dbt.Group,
		EquipmentType: dbt.EquipmentType,
		Type:          dbt.Type,
	}
}
