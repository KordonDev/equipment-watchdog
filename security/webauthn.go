package security

import (
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
var sessinoStore *session.Store

func init() {
	var err error
	webAuthn, err = webauthn.New(&webauthn.Config{
		RPDisplayName: "equipment watchdog", // Display Name for your site
		RPID:          "localhost",          // Generally the domain name for your site
		RPOrigin:      "http://localhost",   // The origin URL for WebAuthn requests
	})

	if err != nil {
		log.Fatal("Error creating webAuthn", err)
	}

	sessinoStore, err = session.NewStore()
	if err != nil {
		log.Fatal("Error creating sessionStore", err)
	}
}

func startRegister(c *gin.Context) {
	var userRequest UserRequest

	if err := c.BindJSON(&userRequest); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	db := getUserDB()

	user, err := db.GetUser(userRequest.username)
	if err != nil {
		user = &User{
			name:        userRequest.username,
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

	c.JSON(http.StatusOK, options)
}

type User struct {
	name        string
	credentials []webauthn.Credential
}

func (u User) ExcludedCredentials() []protocol.CredentialDescriptor {
	excludeList := []protocol.CredentialDescriptor{}

	for _, cred := range u.credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialId: cred.ID,
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
