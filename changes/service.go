package changes

import (
	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

type ChangeDatabase interface {
	save(models.Change) (*models.Change, error)
}

type ChangeService struct {
	db          ChangeDatabase
	userService UserService
}

type UserService interface {
	GetUser(string) (*models.User, error)
}

func NewChangeService(db *gorm.DB, userService UserService) ChangeService {
	return ChangeService{
		db:          newChangeDB(db),
		userService: userService,
	}
}

func (cs ChangeService) Save(change models.Change, c *gin.Context) (*models.Change, error) {
	username := c.GetString("username")
	if user, err := cs.userService.GetUser(username); err != nil {
		change.ByUser = user.ID
	}

	return cs.db.save(change)

}
