package url

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseToString(c *gin.Context, paramName string) (string, error) {
	paramValue := c.Param(paramName)
	if paramValue == "" {
		return "", fmt.Errorf("Param %s not found", paramName)
	}
	return paramValue, nil
}

func ParseToInt(c *gin.Context, paramName string) (uint64, error) {
	paramValueString := c.Param(paramName)
	paramValueNumber, err := strconv.ParseUint(paramValueString, 10, 64)

	if err != nil {
		return 0, fmt.Errorf("%s as number could not be found", paramName)
	}
	return paramValueNumber, nil
}

const authorizationCookieKey = "Authorization"

func SetCookie(c *gin.Context, token string, domain string) {
	c.SetCookie(authorizationCookieKey, token, 60*100, "/", domain, true, true)
}

func RemoveCookie(c *gin.Context, domain string) {
	c.SetCookie(authorizationCookieKey, "", 60*100, "/", domain, true, true)
}
