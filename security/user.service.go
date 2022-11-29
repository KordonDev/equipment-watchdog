package security

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userService struct {
	db         *userDB
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

	fmt.Println(username)
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
