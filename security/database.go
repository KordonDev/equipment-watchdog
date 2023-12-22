package security

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
	"time"
)

type sessionDB struct {
	*gorm.DB
}

func NewDatebase(db *gorm.DB) *sessionDB {
	db.AutoMigrate(&DbSessionData{})

	return &sessionDB{
		db,
	}
}

func (s *sessionDB) getSession(username string) (webauthn.SessionData, error) {
	var dbs DbSessionData
	err := s.Model(DbSessionData{}).First(&dbs, "Username = ?", username).Error

	if err != nil {
		return webauthn.SessionData{}, fmt.Errorf("error getting user: %s", username)
	}

	return dbs.ToSession(), nil
}

func (s *sessionDB) storeSession(username string, sessionData webauthn.SessionData) error {
	dbs := DbSessionData{
		Username:             username,
		Challenge:            sessionData.Challenge,
		UserID:               sessionData.UserID,
		AllowedCredentialIDs: sessionData.AllowedCredentialIDs,
		Expires:              sessionData.Expires,
		UserVerification:     sessionData.UserVerification,
	}

	return s.Save(&dbs).Where("username = ?", username).Error
}


type DbSessionData struct {
	Username             string                               `gorm:"primaryKey"`
	Challenge            string                               `json:"challenge"`
	UserID               []byte                               `json:"user_id"`
	AllowedCredentialIDs ByteArrayArray                       `json:"allowed_credentials,omitempty"`
	Expires              time.Time                            `json:"expires"`
	UserVerification     protocol.UserVerificationRequirement `gorm:"text"`
}

type ByteArrayArray [][]byte
func (b ByteArrayArray) Value() (driver.Value, error) {
	return json.Marshal(b)
}
func (b *ByteArrayArray) Scan(input interface{}) error {
	if input == nil {
		*b = nil
		return nil
	}

	return json.Unmarshal(input.([]byte), b)
}

func (dbs DbSessionData) ToSession() webauthn.SessionData {
	return webauthn.SessionData{
		Challenge:            dbs.Challenge,
		UserID:               dbs.UserID,
		AllowedCredentialIDs: dbs.AllowedCredentialIDs,
		Expires:              dbs.Expires,
		UserVerification:     dbs.UserVerification,
		Extensions:           protocol.AuthenticationExtensions{},
	}
}

func (DbSessionData) TableName() string {
	return "webauthn_sessions"
}
