package security

import (
	"encoding/binary"
	"fmt"
	"log"
	"net/http"

	"github.com/duo-labs/webauthn.io/session"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gin-gonic/gin"
)

type UserRequest struct {
	username string `json:username`
}

var webAuthn *webauthn.WebAuthn
var sessionStore *session.Store

func init() {
	var err error
	webAuthn, err = webauthn.New(&webauthn.Config{
		RPDisplayName: "equipment watchdog",    // Display Name for your site
		RPID:          "localhost",             // Generally the domain name for your site
		RPOrigin:      "http://localhost:8080", // The origin URL for WebAuthn requests
	})

	if err != nil {
		log.Fatal("Error creating webAuthn", err)
	}

	sessionStore, err = session.NewStore()
	if err != nil {
		log.Fatal("Error creating sessionStore", err)
	}
}

func StartRegister(c *gin.Context) {
	username := c.Param("username")

	db := getUserDB()

	user, err := db.GetUser(username)
	if err != nil {
		user = &User{
			name:        username,
			credentials: []webauthn.Credential{},
		}
		db.AddUser(user)
	}

	// TODO: go PublicKeyCredentialCreationOptions
	registerOpts := func(credOptions *protocol.PublicKeyCredentialCreationOptions) {
		credOptions.CredentialExcludeList = user.ExcludedCredentials()
	}

	options, sessionData, err := webAuthn.BeginRegistration(
		user,
		registerOpts,
	)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// store options
	err = sessionStore.SaveWebauthnSession("registration", sessionData, c.Request, c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, options)
}

func FinishRegistration(c *gin.Context) {
	username := c.Param("username")

	user, err := db.GetUser(username)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	sessionData, err := sessionStore.GetWebauthnSession("registration", c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	credential, err := webAuthn.FinishRegistration(user, sessionData, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user.AddCredential(*credential)

	c.Status(http.StatusOK)
}

func StartLogin(c *gin.Context) {
	username := c.Param("username")

	user, err := db.GetUser(username)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	options, sessionData, err := webAuthn.BeginLogin(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = sessionStore.SaveWebauthnSession("authentication", sessionData, c.Request, c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, options)
}

func FinishLogin(c *gin.Context) {
	username := c.Param("username")

	user, err := db.GetUser(username)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	sessionData, err := sessionStore.GetWebauthnSession("authentication", c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	_, err = webAuthn.FinishLogin(user, sessionData, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, "successfull login")
}

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

type userDB struct {
	users map[string]*User
}

var db *userDB

func getUserDB() *userDB {
	if db == nil {
		db = &userDB{
			users: make(map[string]*User),
		}
	}
	return db
}

func (db *userDB) GetUser(username string) (*User, error) {
	user, ok := db.users[username]

	if !ok {
		return &User{}, fmt.Errorf("error getting user: %s", username)
	}

	return user, nil
}

func (db *userDB) AddUser(user *User) {
	db.users[user.name] = user
}
