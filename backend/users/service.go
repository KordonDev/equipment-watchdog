package users

import (
	"context"
	"github.com/kordondev/equipment-watchdog/models"

	"golang.org/x/crypto/bcrypt"
)

type UserDatabase interface {
	getUser(string) (*models.User, error)
	getAll() ([]*models.User, error)
	saveUser(*models.User) (*models.User, error)
	addUser(*models.User) (*models.User, error)
	hasApprovedAndAdminUser() bool
	getForIds([]uint64) ([]*models.User, error)
	changePassword(string, string) error
	getPasswordHashForUser(username string) (string, error)
}

type JwtService interface {
	GenerateToken(models.User) string
}

type userService struct {
	db         UserDatabase
	jwtService JwtService
}

func NewUserService(db *userDB, jwtService JwtService) *userService {
	return &userService{
		db,
		jwtService,
	}
}

func (u *userService) getTokenForUser(user *models.User) string {
	return u.jwtService.GenerateToken(*user)
}

func (u *userService) GetUser(username string) (*models.User, error) {
	return u.db.getUser(username)
}

func (u *userService) getAll() ([]*models.User, error) {
	return u.db.getAll()
}

func (u *userService) toggleApprove(username string) (*models.User, error) {
	user, err := u.db.getUser(username)
	if err != nil {
		return nil, err
	}

	user.IsApproved = !user.IsApproved
	u.db.saveUser(user)
	return user, nil
}

func (u *userService) toggleAdmin(username string) (*models.User, error) {
	user, err := u.db.getUser(username)
	if err != nil {
		return nil, err
	}

	user.IsAdmin = !user.IsAdmin
	u.db.saveUser(user)
	return user, nil
}

func (u *userService) AddUser(user *models.User) (*models.User, error) {
	return u.db.addUser(user)
}

func (u *userService) SaveUser(user *models.User) (*models.User, error) {
	return u.db.saveUser(user)
}

func (u *userService) HasApprovedAndAdminUser() bool {
	return u.db.hasApprovedAndAdminUser()
}

func (u *userService) GetForIds(ids []uint64) ([]*models.User, error) {
	return u.db.getForIds(ids)
}

func (u *userService) ChangePassword(c context.Context, username, password string) error {
	return u.db.changePassword(username, password)
}

func (u *userService) CheckLogin(username, password string) error {
	pwHash, err := u.db.getPasswordHashForUser(username)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(pwHash), []byte(password)); err != nil {
		return err
	}
	return nil
}
