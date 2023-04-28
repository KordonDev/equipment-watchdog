package security

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserDatabase interface {
	GetUser(string) (*User, error)
	GetAll() ([]*User, error)
	SaveUser(*User) (*User, error)
}

// TODO: refactor own user package?
type userService struct {
	db         UserDatabase
	jwtService *JwtService
}

func NewUserService(db *userDB, jwtService *JwtService) *userService {
	return &userService{
		db,
		jwtService,
	}
}

func (u *userService) GetMe(c *gin.Context) {
	username := c.GetString("username")
	user, err := u.db.GetUser(username)

	if err != nil {
		fmt.Printf("%v", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	token := u.jwtService.GenerateToken(*user)
	u.jwtService.SetCookie(c, token)
	c.JSON(http.StatusOK, user)
}

func (u *userService) GetAll(c *gin.Context) {
	users, err := u.db.GetAll()

	if err != nil {
		fmt.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

func (u *userService) ToggleApprove(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := u.db.GetUser(username)
	if err != nil {
		fmt.Println(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	user.IsApproved = !user.IsApproved
	u.db.SaveUser(user)
	c.JSON(http.StatusOK, user)
}

func (u *userService) ToggleAdmin(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := u.db.GetUser(username)
	if err != nil {
		fmt.Println(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	user.IsAdmin = !user.IsAdmin
	u.db.SaveUser(user)
	c.JSON(http.StatusOK, user)
}
