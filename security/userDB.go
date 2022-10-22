package security

import "fmt"

type userDB struct {
	users map[string]*User
}

func NewUserDB() *userDB {
	return &userDB{
		users: make(map[string]*User),
	}
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
