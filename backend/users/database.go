package users

import (
	"fmt"

	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

type userDB struct {
	*gorm.DB
}

func NewDatebase(db *gorm.DB) *userDB {
	return &userDB{
		db,
	}
}

func (u *userDB) getUser(username string) (*models.User, error) {
	var dbu models.DbUser
	err := u.Model(&models.DbUser{}).Preload("Credentials").First(&dbu, "name = ?", username).Error

	if err != nil {
		return &models.User{}, fmt.Errorf("error getting user: %s", username)
	}

	return dbu.ToUser(), nil
}

func (u *userDB) getAll() ([]*models.User, error) {
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

func (u *userDB) addUser(user *models.User) (*models.User, error) {
	us := user.ToDBUser()
	err := u.Create(us).Error
	if err != nil {
		return nil, err
	}
	return us.ToUser(), nil
}

func (u *userDB) saveUser(user *models.User) (*models.User, error) {
	u.Save(user.ToDBUser())
	return u.getUser(user.Name)
}

func (u *userDB) hasApprovedAndAdminUser() bool {
	var dbu models.DbUser
	err := u.Model(&models.DbUser{}).First(&dbu, "is_admin = 1 AND is_approved = 1").Error
	if err != nil || dbu.ID == 0 {
		return false
	}
	return true
}

func (u *userDB) getForIds(ids []uint64) ([]*models.User, error) {
	dbUser := make([]*models.DbUser, 0)

	err := u.Where("ID IN ?", ids).Find(&dbUser).Error
	if err != nil {
		return make([]*models.User, 0), err
	}

	return listFromDb(dbUser), err
}

func listFromDb(dbu []*models.DbUser) []*models.User {
	users := make([]*models.User, 0)
	for _, user := range dbu {
		users = append(users, user.ToUser())
	}
	return users
}

func (u *userDB) changePassword(username, password string) error {
	return u.DB.Exec("UPDATE users SET password = ? WHERE name = ?", password, username).Error
}

type password struct {
	Password string `gorm:"password"`
}

func (password) TableName() string {
	return "users"
}

func (u *userDB) getPasswordHashForUser(username string) (string, error) {
	var p password
	err := u.Where("name = ?", username).First(&p).Error
	if err != nil {
		return "", err
	}

	return p.Password, nil
}

func (u *userDB) getUserByCredentialId(credentialId string) (*models.User, error) {
	var dbu models.DbUser
	err := u.Model(&models.DbUser{}).Preload("Credentials").First(&dbu, "credential_id = ?", credentialId).Error
	if err != nil {
		return &models.User{}, fmt.Errorf("error getting user by credential_id: %w", err)
	}
	return dbu.ToUser(), nil
}
