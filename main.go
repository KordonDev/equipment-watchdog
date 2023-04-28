package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/config"
	"github.com/kordondev/equipment-watchdog/equipment"
	"github.com/kordondev/equipment-watchdog/members"
	"github.com/kordondev/equipment-watchdog/security"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	configuration, err := config.New("./configuration")
	if err != nil {
		panic(err)
	}

	db, err := createDB(configuration.Debug, configuration.DSN)
	if err != nil {
		panic(err)
	}

	userDB := security.NewUserDB(db)

	jwtService := security.NewJwtService(configuration.Origin, configuration.JwtSecret, configuration.Domain)

	//TODO: discussion about 'who delivers the frontend?' - Multiple options
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{configuration.Origin}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	api := router.Group("/api")

	//TODO: webAuthController for all of the route mapping
	webAuthNService, err := security.NewWebAuthNService(userDB, configuration.Origin, configuration.Domain, jwtService)
	if err != nil {
		panic(fmt.Sprintf("Error creating webAuthn: %v", err))
	}
	api.GET("/register/:username", webAuthNService.StartRegister)
	api.POST("/register/:username", webAuthNService.FinishRegistration)
	api.GET("/login/:username", webAuthNService.StartLogin)
	api.POST("/login/:username", webAuthNService.FinishLogin)
	api.POST("/logout", webAuthNService.Logout)
	api.Use(security.AuthorizeJWTMiddleware(configuration.Origin, jwtService))

	equipmentService := equipment.NewEquipmentService(db)
	database := members.NewMemberDB(db)
	memberService := members.NewMemberService(database, equipmentService)
  membersRoute := api.Group("/members")
	{
		membersRoute.GET("/", memberService.GetAllMembers)
		membersRoute.POST("/", memberService.CreateMember)

		membersRoute.GET("/:id", memberService.GetMemberById)
		membersRoute.PUT("/:id", memberService.UpdateMember)
		membersRoute.DELETE("/:id", memberService.DeleteById)

		membersRoute.GET("/groups", memberService.GetAllGroups)
	}

	userService := security.NewUserService(userDB, jwtService)
	api.GET("/me", userService.GetMe)
	userRoute := api.Group("/users")
	{
		userRoute.GET("/", security.AdminOnlyMiddleware(), userService.GetAll)

		userRoute.PATCH("/:username/toggle-approve", security.AdminOnlyMiddleware(), userService.ToggleApprove)
		userRoute.PATCH("/:username/toggle-admin", security.AdminOnlyMiddleware(), userService.ToggleAdmin)
	}

	equipmentRoute := api.Group("/equipment")
	{
		equipmentRoute.POST("/", equipmentService.CreateEquipment)

		equipmentRoute.GET("/:id", equipmentService.GetEquipmentById)
		equipmentRoute.DELETE("/:id", equipmentService.DeleteEquipment)

		equipmentRoute.GET("/free", equipmentService.FreeEquipment)

		equipmentRoute.GET("/type/:type", equipmentService.GetAllEquipmentByType)
	}

	router.Run(fmt.Sprintf("%s:8080", configuration.Domain))
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
