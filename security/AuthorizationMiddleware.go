package security

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthorizeJWTMiddleware(domain string, jwtService *JwtService) gin.HandlerFunc {
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
			token, err := jwtService.ValidateToken(jwtCookie)
			if token.Valid {
				jwtData := jwtService.GetClaims(token)

				newToken := jwtService.GenerateToken(jwtData.Name, jwtData.IsUser)
				c.SetCookie(AUTHORIZATION_COOKIE_KEY, newToken, 60*100, "/", domain, true, true)
			} else {
				fmt.Println(err)

				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"redirect": "login",
				})
			}
		} else {
			fmt.Println("No auth header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"redirect": "login",
			})
		}
	}
}
