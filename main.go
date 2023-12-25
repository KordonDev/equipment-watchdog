package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/kordondev/equipment-watchdog/changes"
	"github.com/kordondev/equipment-watchdog/config"
	"github.com/kordondev/equipment-watchdog/equipment"
	"github.com/kordondev/equipment-watchdog/members"
	"github.com/kordondev/equipment-watchdog/orders"
	"github.com/kordondev/equipment-watchdog/registrationcodes"
	"github.com/kordondev/equipment-watchdog/security"
	"github.com/kordondev/equipment-watchdog/users"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	configuration, err := config.New("./configuration")
	if err != nil {
		panic(err)
	}

	db, err := createDB(configuration.Debug, configuration.DatabaseConnection)
	if err != nil {
		panic(err)
	}

	userDB := users.NewDatebase(db)

	jwtService := security.NewJwtService(configuration.Origin, configuration.JwtSecret)

	router := gin.Default()
	router.SetTrustedProxies(nil)

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{configuration.Origin}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	api := router.Group("/api")

	api.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	userService := users.NewUserService(userDB, jwtService)
	webAuthNService, err := security.NewWebAuthNService(userService, configuration.Origin, configuration.Domain, jwtService, db)
	if err != nil {
		panic(fmt.Sprintf("Error creating webAuthn: %v", err))
	}
	security.NewController(api, webAuthNService, configuration.Domain)
	// TODO: security.Controller
	api.Use(security.AuthorizeJWTMiddleware(configuration.Origin, jwtService, configuration.Domain))

	changeWriter := changes.NewChangeWriterService(db, userService)

	equipmentService := equipment.NewEquipmentService(db)
	database := members.NewMemberDB(db)
	memberService := members.NewMemberService(database, &equipmentService)
	members.NewController(api, memberService, changeWriter)

	users.NewController(api, userService, configuration.Domain)
	equipment.NewController(api, equipmentService, changeWriter)

	orderService := orders.NewOrderService(db, &equipmentService)
	orders.NewController(api, orderService, changeWriter)

	registrationCodesService := registrationcodes.NewService(db, equipmentService)
	registrationcodes.NewController(api, registrationCodesService)

	changeService := changes.NewChangeService(db, equipmentService, memberService, userService, orderService)
	changes.NewController(api, changeService)

	router.Run(":8080")
}

func createDB(debug bool, dsn string) (*gorm.DB, error) {
	logLevel := logger.Error
	if debug {
		logLevel = logger.Info
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logLevel,    // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, errors.New("failed to connect database")
	}
	return db, nil
}
