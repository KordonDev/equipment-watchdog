package security

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type userDB struct {
	db *gorm.DB
}

func NewUserDB() *userDB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{})

	return &userDB{
		db: db,
	}
}

func (u *userDB) GetUser(username string) (*User, error) {
	var user User
	err := u.db.First(&user, "name = ?", username).Error

	if err != nil {
		return &User{}, fmt.Errorf("error getting user: %s", username)
	}

	return &user, nil
}

func (u *userDB) AddUser(user *User) {
	u.db.Create(user)
}
