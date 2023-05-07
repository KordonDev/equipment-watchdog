package security

import (
	"fmt"

	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

// TODO: refactor own user package?
type userDB struct {
	*gorm.DB
}

func NewUserDB(db *gorm.DB) *userDB {
	db.AutoMigrate(&models.DbUser{}, &models.DbCredential{}, &models.DbAuthenticator{})

	return &userDB{
		db,
	}
}

func (u *userDB) GetUser(username string) (*models.User, error) {
	var dbu models.DbUser
	err := u.Model(&models.DbUser{}).Preload("Credentials").First(&dbu, "name = ?", username).Error

	if err != nil {
		return &models.User{}, fmt.Errorf("error getting user: %s", username)
	}

	return dbu.ToUser(), nil
}

func (u *userDB) GetAll() ([]*models.User, error) {
	var dbUser []models.DbUser
	err := u.Find(&dbUser).Error

	if err != nil {
		return nil, err
	}

	users := make([]*models.User, 0)
	for _, user := range dbUser {
		users = append(users, user.ToUser())
	}
	return users, nil
}

func (u *userDB) AddUser(user *models.User) (*models.User, error) {
	us := user.ToDBUser()
	err := u.Create(us).Error
	if err != nil {
		return nil, err
	}
	return us.ToUser(), nil
}

func (u *userDB) SaveUser(user *models.User) (*models.User, error) {
	u.Save(user.ToDBUser())
	return u.GetUser(user.Name)
}

func (u *userDB) HasApprovedAndAdminUser() bool {
	var dbu models.DbUser
	err := u.Model(&models.DbUser{}).First(&dbu, "is_admin = 1 AND is_approved = 1").Error
	if err != nil || dbu.ID == 0 {
		return false
	}
	return true
}
