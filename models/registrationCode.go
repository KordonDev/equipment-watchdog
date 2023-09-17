package models

import "time"

type RegistrationCode struct {
	ID            string    `json:"id"`
	ReservedUntil time.Time `json:"reservedUntil"`
}

type DbRegistrationCode struct {
	ID            string `gorm:"primarykey"`
	ReservedUntil time.Time
}

func (DbRegistrationCode) TableName() string {
	return "registration_codes"
}

func (r *RegistrationCode) ToDb() *DbRegistrationCode {
	return &DbRegistrationCode{
		ID:            r.ID,
		ReservedUntil: r.ReservedUntil,
	}
}

func (r *DbRegistrationCode) FromDb() *RegistrationCode {
	return &RegistrationCode{
		ID:            r.ID,
		ReservedUntil: r.ReservedUntil,
	}
}
