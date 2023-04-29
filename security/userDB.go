package security

import (
	"fmt"
	"time"

	"github.com/duo-labs/webauthn/webauthn"
	"gorm.io/gorm"
)

// TODO: refactor own user package?
type userDB struct {
	*gorm.DB
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

type DbUser struct {
	ID          uint64 `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string `gorm:"unique"`
	IsApproved  bool
	IsAdmin     bool
	Credentials []DbCredential `gorm:"foreignKey:UserID"`
}

func (DbUser) TableName() string {
	return "users"
}

func NewUserDB(db *gorm.DB) *userDB {
	db.AutoMigrate(&DbUser{}, &DbCredential{}, &DbAuthenticator{})

	return &userDB{
		db,
	}
}

func (u *userDB) GetUser(username string) (*User, error) {
	var dbu DbUser
	err := u.Model(&DbUser{}).Preload("Credentials").First(&dbu, "name = ?", username).Error

	if err != nil {
		return &User{}, fmt.Errorf("error getting user: %s", username)
	}

	return dbu.toUser(), nil
}

func (u *userDB) GetAll() ([]*User, error) {
	var dbUser []DbUser
	err := u.Find(&dbUser).Error

	if err != nil {
		return nil, err
	}

	users := make([]*User, 0)
	for _, user := range dbUser {
		users = append(users, user.toUser())
	}
	return users, nil
}

func (u *userDB) AddUser(user *User) (*User, error) {
	us := user.toDBUser()
	err := u.Create(us).Error
	if err != nil {
		return nil, err
	}
	return us.toUser(), nil
}

func (u *userDB) SaveUser(user *User) (*User, error) {
	u.Save(user.toDBUser())
	return u.GetUser(user.Name)
}

func (u *userDB) HasApprovedAndAdminUser() bool {
	var dbu DbUser
	err := u.Model(&DbUser{}).First(&dbu, "is_admin = 1 AND is_approved = 1").Error
	if err != nil || dbu.ID == 0 {
		return false
	}
	return true
}

func (u *User) toDBUser() *DbUser {
	var c []DbCredential
	for _, cr := range u.Credentials {
		a := DbAuthenticator{
			AAGUID:       cr.Authenticator.AAGUID,
			SignCount:    cr.Authenticator.SignCount,
			CloneWarning: cr.Authenticator.CloneWarning,
		}

		dbC := DbCredential{
			ID:              cr.ID,
			UserID:          u.ID,
			PublicKey:       cr.PublicKey,
			AttestationType: cr.AttestationType,
			Authenticator:   a,
		}
		c = append(c, dbC)
	}
	dbu := DbUser{
		ID:          u.ID,
		Name:        u.Name,
		IsApproved:  u.IsApproved,
		IsAdmin:     u.IsAdmin,
		Credentials: c,
	}
	return &dbu
}

func (dbu *DbUser) toUser() *User {
	var c []webauthn.Credential
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
