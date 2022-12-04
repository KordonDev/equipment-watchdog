package security

import (
	"encoding/binary"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
)

type User struct {
	ID          uint64                `json:"id"`
	Name        string                `json:"name" mapstructure:"name"`
	IsApproved  bool                  `json:"isApproved"`
	IsAdmin     bool                  `json:"isAdmin"`
	Credentials []webauthn.Credential `json:"-"`
}

func NewUser(name string) *User {
	return &User{
		Name:        name,
		Credentials: []webauthn.Credential{},
	}
}

func (u User) WebAuthnID() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(u.ID))
	binary.LittleEndian.Uint64(buf)
	return buf
}

func (u User) WebAuthnName() string {
	return u.Name
}

func (u User) WebAuthnDisplayName() string {
	return u.Name
}

func (u User) WebAuthnIcon() string {
	return ""
}

func (u User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

func (u *User) AddCredential(c webauthn.Credential) {
	u.Credentials = append(u.Credentials, c)
}

func (u User) ExcludedCredentials() []protocol.CredentialDescriptor {
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
