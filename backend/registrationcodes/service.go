package registrationcodes

import (
	"math/rand"
	"time"

	"github.com/cloudflare/cfssl/log"
	"github.com/go-co-op/gocron"

	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

type database interface {
	save(models.RegistrationCode) error
	exists(string) bool
	deleteOutdated() error
}

type equipmentService interface {
	RegistrationCodeExists(string) bool
}

type Service struct {
	db               database
	equipmentService equipmentService
}

func NewService(db *gorm.DB, equipmentService equipmentService) Service {
	registrationCodesDB := newDatabase(db)

	rcs := Service{
		db:               registrationCodesDB,
		equipmentService: equipmentService,
	}

	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Hour().Do(rcs.deleteOutdated())
	s.StartAsync()

	return rcs
}

func (s Service) getRegistrationCode() (models.RegistrationCode, error) {
	registrationCode := s.createRandomRegistrationCode()

	if err := s.save(registrationCode); err != nil {
		return models.RegistrationCode{}, err
	}

	return registrationCode, nil
}

func (s Service) save(registrationCode models.RegistrationCode) error {
	return s.db.save(registrationCode)
}

func (s Service) deleteOutdated() error {
	log.Info("delete outdated rc (Now: ", time.Now().Format("2006-01-02 15:04:05"), ")")
	return s.db.deleteOutdated()
}

func (s Service) createRandomRegistrationCode() models.RegistrationCode {
	exists := true
	var ID string
	for exists {
		ID = randomString(4)
		exists = s.db.exists(ID) || s.equipmentService.RegistrationCodeExists(ID)
	}
	return models.RegistrationCode{
		ID:            ID,
		ReservedUntil: time.Now().Add(time.Hour),
	}
}

const firstLetter = "abcdefghiklmnopqrstuwxyz123456789123456789"
const allLetters = "0" + firstLetter

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		if i == 0 {
			b[i] = firstLetter[rand.Intn(len(firstLetter))]
		} else {
			b[i] = allLetters[rand.Intn(len(allLetters))]
		}
	}
	return string(b)
}
