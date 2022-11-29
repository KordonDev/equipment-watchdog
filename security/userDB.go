package security

import (
	"fmt"
	"time"

	"github.com/duo-labs/webauthn/webauthn"
	"gorm.io/gorm"
)

type userDB struct {
	db *gorm.DB
}

type DbAuthenticator struct {
	AAGUID       []byte
	SignCount    uint32
	CloneWarning bool
}

func (DbAuthenticator) TableName() string {
	return "user_credential_authenticators"
}

type DbCredential struct {
	ID              []byte `gorm:"primarykey"`
	UserID          uint64
	PublicKey       []byte
	AttestationType string
	Authenticator   DbAuthenticator `gorm:"embedded;embeddedPrefix:authenticator_"`
	CreatedAt       time.Time
}

func (DbCredential) TableName() string {
	return "user_credentials"
}

type dbUser struct {
	ID          uint64 `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string `gorm:"unique"`
	IsApproved  bool
	IsAdmin     bool
	Credentials []DbCredential `gorm:"foreignKey:UserID"`
}

func (dbUser) TableName() string {
	return "users"
}

func NewUserDB(db *gorm.DB) *userDB {
	db.AutoMigrate(&dbUser{}, &DbCredential{}, &DbAuthenticator{})

	return &userDB{
		db: db,
	}
}

func (u *userDB) GetUser(username string) (*User, error) {
	var dbu dbUser
	err := u.db.Model(&dbUser{}).Preload("Credentials").First(&dbu, "name = ?", username).Error

	if err != nil {
		return &User{}, fmt.Errorf("error getting user: %s", username)
	}

	return dbu.toUser(), nil
}

func (u *userDB) GetAll() ([]*User, error) {
	var dbUser []dbUser
	err := u.db.Find(&dbUser).Error

	if err != nil {
		return nil, err
	}

	var users []*User
	for _, user := range dbUser {
		users = append(users, user.toUser())
	}
	return users, nil
}

func (u *userDB) AddUser(user *User) (*User, error) {
	u.db.Create(user.toDBUser())
	return u.GetUser(user.Name)
}

func (u *userDB) SaveUser(user *User) {
	u.db.Save(user.toDBUser())
}

func (user User) toDBUser() dbUser {
	c := []DbCredential{}
	for _, cr := range user.Credentials {
		a := DbAuthenticator{
			AAGUID:       cr.Authenticator.AAGUID,
			SignCount:    cr.Authenticator.SignCount,
			CloneWarning: cr.Authenticator.CloneWarning,
		}

		dbC := DbCredential{
			ID:              cr.ID,
			UserID:          user.ID,
			PublicKey:       cr.PublicKey,
			AttestationType: cr.AttestationType,
			Authenticator:   a,
		}
		c = append(c, dbC)
	}
	dbu := dbUser{
		ID:          user.ID,
		Name:        user.Name,
		IsApproved:  user.IsApproved,
		IsAdmin:     user.IsAdmin,
		Credentials: c,
	}
	return dbu
}

func (dbu dbUser) toUser() *User {
	c := []webauthn.Credential{}
	for _, cr := range dbu.Credentials {
		a := webauthn.Authenticator{
			AAGUID:       cr.Authenticator.AAGUID,
			SignCount:    cr.Authenticator.SignCount,
			CloneWarning: cr.Authenticator.CloneWarning,
		}

		dbC := webauthn.Credential{
			ID:              cr.ID,
			PublicKey:       cr.PublicKey,
			AttestationType: cr.AttestationType,
			Authenticator:   a,
		}
		c = append(c, dbC)
	}
	user := User{
		ID:          dbu.ID,
		Name:        dbu.Name,
		IsApproved:  dbu.IsApproved,
		IsAdmin:     dbu.IsAdmin,
		Credentials: c,
	}
	return &user
}
