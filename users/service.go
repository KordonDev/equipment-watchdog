package users

import (
	"github.com/kordondev/equipment-watchdog/models"
)

type UserDatabase interface {
	GetUser(string) (*models.User, error)
	GetAll() ([]*models.User, error)
	SaveUser(*models.User) (*models.User, error)
	AddUser(*models.User) (*models.User, error)
	HasApprovedAndAdminUser() bool
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

func (u *userService) GetUserWithToken(username string) (*models.User, string, error) {
	user, err := u.GetUser(username)
	if err != nil {
		return nil, "", err
	}
	token := u.jwtService.GenerateToken(*user)
	return user, token, nil
}

func (u *userService) GetUser(username string) (*models.User, error) {
	return u.db.GetUser(username)
}

func (u *userService) GetAll() ([]*models.User, error) {
	return u.db.GetAll()
}

func (u *userService) ToggleApprove(username string) (*models.User, error) {
	user, err := u.db.GetUser(username)
	if err != nil {
		return nil, err
	}

	user.IsApproved = !user.IsApproved
	u.db.SaveUser(user)
	return user, nil
}

func (u *userService) ToggleAdmin(username string) (*models.User, error) {
	user, err := u.db.GetUser(username)
	if err != nil {
		return nil, err
	}

	user.IsAdmin = !user.IsAdmin
	u.db.SaveUser(user)
	return user, nil
}

func (u *userService) AddUser(user *models.User) (*models.User, error) {
	return u.db.AddUser(user)
}

func (u *userService) SaveUser(user *models.User) (*models.User, error) {
	return u.db.SaveUser(user)
}

func (u *userService) HasApprovedAndAdminUser() bool {
	return u.db.HasApprovedAndAdminUser()
}
