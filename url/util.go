package url

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseToInt(c *gin.Context, paramName string) (uint64, error) {
	paramValueString := c.Param(paramName)
	paramValueNumber, err := strconv.ParseUint(paramValueString, 10, 64)

	if err != nil {
		return 0, fmt.Errorf("%s as number could not be found", paramName)
	}
	return paramValueNumber, nil
}
