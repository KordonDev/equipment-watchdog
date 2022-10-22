package security

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthorizeJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		var jwtCookie string
		cookie := c.GetHeader("Cookie")
		cookies := strings.Split(cookie, "; ")
		for _, c := range cookies {
			if strings.HasPrefix(c, "Authorization") {
				jwtCookie = strings.Split(c, "=")[1]
			}
		}

		if len(jwtCookie) > 0 {
			token, err := JWTAuthService().ValidateToken(jwtCookie)
			if token.Valid {
				claims := token.Claims.(jwt.MapClaims)
				fmt.Println(claims)
			} else {
				fmt.Println(err)
				c.AbortWithStatus(http.StatusUnauthorized)
			}
		} else {
			fmt.Println("No auth header")
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
