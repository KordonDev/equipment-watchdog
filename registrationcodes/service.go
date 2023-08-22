package registrationcodes

import (
	"time"

	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

type database interface {
  save(models.RegistrationCode) error
}

type service struct {
  db database
}

func NewService(db *gorm.DB) *service {
  registrationCodesDB := newDatabase(db)

  return &service{
    db: registrationCodesDB,
  }
}

func (s service) getRegistrationCode() (models.RegistrationCode, error) {
  registrationCode := createRandomRegistrationCode()

  if err := s.save(registrationCode); err != nil {
    return models.RegistrationCode{}, err
  }

  return registrationCode, nil
}

func (s service) save(registrationCode models.RegistrationCode) error {
  return s.db.save(registrationCode);
}

func createRandomRegistrationCode() models.RegistrationCode {
  // TODO: Loop if code exists in registrationCodesDB or equipmentDB
  return models.RegistrationCode{
    ID: "123",
    ReservedUntil: time.Now().Add(time.Hour),
  }

}
