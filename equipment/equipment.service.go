package equipment

import (
	"errors"
	"github.com/cloudflare/cfssl/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type EquipmentService struct {
	db *equipmentDB
}

func NewEquipmentService(equipmentDB *equipmentDB) *EquipmentService {
	return &EquipmentService{
		db: equipmentDB,
	}
}

func (e *EquipmentService) GetEquipmentById(c *gin.Context) {
	id, err := parseId(c)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}


	// TODO: go

}

func parseId(c *gin.Context) (uint64, error) {
	id := c.Param("id")
	idN, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		return 0, errors.New("id as number could not be found")
	}
	return idN, nil
}


