package models

import (
	"encoding/binary"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

type User struct {
	ID            uint64                `json:"id"`
	Name          string                `json:"name" mapstructure:"name"`
	IsApproved    bool                  `json:"isApproved"`
	IsAdmin       bool                  `json:"isAdmin"`
	Credentials   []webauthn.Credential `json:"-"`
	CredentialIds []string
}

func (u *User) WebAuthnID() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(u.ID))
	binary.LittleEndian.Uint64(buf)
	return buf
}

func (u *User) WebAuthnName() string {
	return u.Name
}

func (u *User) WebAuthnDisplayName() string {
	return u.Name
}

func (u *User) WebAuthnIcon() string {
	return ""
}

func (u *User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

func (u *User) AddCredential(c webauthn.Credential, credentialId string) {
	u.Credentials = append(u.Credentials, c)
	u.CredentialIds = []string{credentialId} // append(u.CredentialIds, credentialId)
}

func (u *User) ExcludedCredentials() []protocol.CredentialDescriptor {
	excludeList := []protocol.CredentialDescriptor{}

	for _, cred := range u.Credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		excludeList = append(excludeList, descriptor)
	}
	return excludeList
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

func (DbCredential) TableName() string {
	return "user_credentials"
}

type DbUser struct {
	ID            uint64 `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string `gorm:"unique"`
	Password      string `gorm:"default:''"`
	IsApproved    bool
	IsAdmin       bool
	Credentials   []DbCredential `gorm:"foreignKey:UserID"`
	CredentialIds []CredentialId `gorm:"foreignKey:UserId"`
}

type CredentialId struct {
	UserId       uint64 `gorm:"primaryKey;autoIncrement:false"`
	CredentialId string `gorm:"primaryKey;autoIncrement:false"`
}

func (DbUser) TableName() string {
	return "users"
}

func (u *User) ToDBUser() *DbUser {
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

	credentialIds := make([]CredentialId, len(u.CredentialIds))
	for i, cred := range u.CredentialIds {
		credentialIds[i] = CredentialId{
			UserId:       u.ID,
			CredentialId: cred,
		}
	}

	dbu := DbUser{
		ID:            u.ID,
		Name:          u.Name,
		IsApproved:    u.IsApproved,
		IsAdmin:       u.IsAdmin,
		Credentials:   c,
		CredentialIds: credentialIds,
	}
	return &dbu
}

func (dbu *DbUser) toSmallUser() *User {
	user := User{
		ID:          dbu.ID,
		Name:        dbu.Name,
		IsApproved:  dbu.IsApproved,
		IsAdmin:     dbu.IsAdmin,
		Credentials: nil,
	}
	return &user
}

func (dbu *DbUser) ToUser() *User {
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

	credentialIds := make([]string, len(dbu.CredentialIds))
	for i, cred := range dbu.CredentialIds {
		credentialIds[i] = cred.CredentialId
	}

	user := User{
		ID:            dbu.ID,
		Name:          dbu.Name,
		IsApproved:    dbu.IsApproved,
		IsAdmin:       dbu.IsAdmin,
		Credentials:   c,
		CredentialIds: credentialIds,
	}
	return &user
}
