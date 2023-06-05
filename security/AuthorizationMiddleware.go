package security

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/url"
)

func AuthorizeJWTMiddleware(_ string, jwtService *JwtService, domain string) gin.HandlerFunc {
	return func(c *gin.Context) {

		var jwtCookie string
		cookie := c.GetHeader("Cookie")
		cookies := strings.Split(cookie, "; ")
		for _, c := range cookies {
			if strings.HasPrefix(c, "Authorization") {
				jwtCookie = strings.Split(c, "=")[1]
			}
		}

		if len(jwtCookie) == 0 {
			fmt.Println("No auth header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"redirect": "login",
			})
			return
		}

		token, err := jwtService.validateToken(jwtCookie)
		if !token.Valid || err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"redirect": "login",
			})
			return
		}

		jwtData, err := jwtService.getClaims(token)
		if !jwtData.IsApproved {
			fmt.Println("Token not approved")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"redirect": "not-approved",
			})
			return
		}
		if err != nil {
			fmt.Println("Error parsing token data", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		newToken := jwtService.GenerateToken(jwtData.User)
		c.Set("username", jwtData.Name)
		c.Set("isApproved", jwtData.IsApproved)
		c.Set("isAdmin", jwtData.IsAdmin)
		url.SetCookie(c, newToken, domain)
	}
}
