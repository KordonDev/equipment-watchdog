package users

import (
	"github.com/kordondev/equipment-watchdog/models"
)

type UserDatabase interface {
	getUser(string) (*models.User, error)
	getAll() ([]*models.User, error)
	saveUser(*models.User) (*models.User, error)
	addUser(*models.User) (*models.User, error)
	hasApprovedAndAdminUser() bool
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
