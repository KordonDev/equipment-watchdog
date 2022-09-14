package security

import (
	"encoding/binary"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
)

type User struct {
	id          uint64
	name        string
	credentials []webauthn.Credential
}

func (u User) WebAuthnID() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(u.id))
	binary.LittleEndian.Uint64(buf)
	return buf
}

func (u User) WebAuthnName() string {
	return u.name
}

func (u User) WebAuthnDisplayName() string {
	return u.name
}

func (u User) WebAuthnIcon() string {
	return ""
}

func (u User) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

func (u *User) AddCredential(c webauthn.Credential) {
	u.credentials = append(u.credentials, c)
}

func (u User) ExcludedCredentials() []protocol.CredentialDescriptor {
	excludeList := []protocol.CredentialDescriptor{}

	for _, cred := range u.credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		excludeList = append(excludeList, descriptor)
	}
	return excludeList
}
