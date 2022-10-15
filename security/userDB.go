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

type DbCredential struct {
	ID              []byte `gorm:"primarykey"`
	UserID          uint64
	PublicKey       []byte
	AttestationType string
	Authenticator   DbAuthenticator `gorm:"embedded;embeddedPrefix:authenticator_"`
	CreatedAt       time.Time
}

type DbUser struct {
	ID          uint64 `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string         `gorm:"unique"`
	Credentials []DbCredential `gorm:"foreignKey:UserID"`
}

func NewUserDB(db *gorm.DB) *userDB {
	db.AutoMigrate(&DbUser{}, &DbCredential{}, &DbAuthenticator{})

	return &userDB{
		db: db,
	}
}

func (u *userDB) GetUser(username string) (*User, error) {
	var dbu DbUser
	err := u.db.Model(&DbUser{}).Preload("Credentials").First(&dbu, "name = ?", username).Error

	if err != nil {
		return &User{}, fmt.Errorf("error getting user: %s", username)
	}

	return fromDBUser(dbu), nil
}

func (u *userDB) AddUser(user *User) (*User, error) {
	u.db.Create(toDBUser(user))
	return u.GetUser(user.Name)
}

func (u *userDB) SaveUser(user *User) {
	u.db.Save(toDBUser(user))
}

func toDBUser(user *User) *DbUser {
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
	dbu := DbUser{
		ID:          user.ID,
		Name:        user.Name,
		Credentials: c,
	}
	return &dbu
}

func fromDBUser(dbu DbUser) *User {
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
		Credentials: c,
	}
	return &user
}