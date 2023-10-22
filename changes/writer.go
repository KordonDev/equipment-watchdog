package changes

import (
	"github.com/cloudflare/cfssl/log"
	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/models"
	"gorm.io/gorm"
)

type ChangeWriterDatabase interface {
	save(models.Change) (*models.Change, error)
}

type ChangeWriterService struct {
	db          ChangeWriterDatabase
	userService UserService
}

type UserService interface {
	GetUser(string) (*models.User, error)
}

func NewChangeWriterService(db *gorm.DB, userService UserService) ChangeWriterService {
	return ChangeWriterService{
		db:          newChangeDB(db),
		userService: userService,
	}
}

func (cs ChangeWriterService) Save(change models.Change, c *gin.Context) (*models.Change, error) {
	username := c.GetString("username")
	if user, err := cs.userService.GetUser(username); err == nil {
		change.UserId = user.ID
	}

	ch, err := cs.db.save(change)
	log.Infof("Changed %+v\n", ch)

	return ch, err
}
