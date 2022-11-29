package security

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		isAdmin := c.GetBool("isAdmin")
		if !isAdmin {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
}
