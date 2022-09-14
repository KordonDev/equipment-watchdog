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
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")

		authHeader2 := c.GetHeader("Cookie")
		cookies := strings.Split(authHeader2, "; ")
		for _, c := range cookies {
			if strings.HasPrefix(c, "Authorization2") {
				token := strings.Split(c, "=")[1]
				fmt.Println(token)
			}
		}

		if len(authHeader) > len(BEARER_SCHEMA) {
			tokenString := authHeader[len(BEARER_SCHEMA):]
			token, err := JWTAuthService().ValidateToken(tokenString)
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
