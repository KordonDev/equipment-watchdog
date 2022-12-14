package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/equipment"
	"github.com/kordondev/equipment-watchdog/members"
	"github.com/kordondev/equipment-watchdog/security"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	config := parseConfig()

	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{config.Origin}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	api := router.Group("/api")

	db := createDB(config.Debug)

	jwtService := security.NewJwtService(config.Origin, config.JwtSecret, config.Domain)
	userDB := security.NewUserDB(db)
	webAuthNService := security.NewWebAuthNService(userDB, config.Origin, config.Domain, jwtService)
	api.GET("/register/:username", webAuthNService.StartRegister)
	api.POST("/register/:username", webAuthNService.FinishRegistration)
	api.GET("/login/:username", webAuthNService.StartLogin)
	api.POST("/login/:username", webAuthNService.FinishLogin)
	api.POST("/logout", webAuthNService.Logout)

	api.Use(security.AuthorizeJWTMiddleware(config.Origin, jwtService))

	memberDB := members.NewMemberDB(db)
	memberService := members.NewMemberService(memberDB)
	membersRoute := api.Group("/members")

	membersRoute.GET("/groups", memberService.GetAllGroups)

	membersRoute.GET("/", memberService.GetAllMembers)
	membersRoute.GET("/:id", memberService.GetMemberById)
	membersRoute.POST("/", memberService.CreateMember)
	membersRoute.PUT("/:id", memberService.UpdateMember)
	membersRoute.DELETE("/:id", memberService.DeleteById)

	userService := security.NewUserService(userDB, jwtService)
	api.GET("/me", userService.GetMe)
	api.PATCH("/users/:username/toggle-approve", security.AdminOnlyMiddleware(), userService.ToggleApprove)
	api.PATCH("/users/:username/toggle-admin", security.AdminOnlyMiddleware(), userService.ToggleAdmin)
	api.GET("/users/", security.AdminOnlyMiddleware(), userService.GetAll)

	equipmentDB := equipment.NewEquipmentDB(db)
	equipmentService := equipment.NewEquipmentService(equipmentDB)
	equipmentRoute := api.Group("/equipment")
	equipmentRoute.GET("/type/:type", equipmentService.GetAllEquipmentByType)
	equipmentRoute.GET("/:id", equipmentService.GetEquipmentById)
	equipmentRoute.POST("/", equipmentService.CreateEquipment)

	router.Run(fmt.Sprintf("%s:8080", config.Domain))
}

func createDB(debug bool) *gorm.DB {
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
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

type Config struct {
	Debug     bool
	Domain    string
	Origin    string
	JwtSecret string
}

func parseConfig() *Config {
	configFile, err := os.ReadFile("./config.yml")
	if err != nil {
		log.Fatal("failed to read config")
	}
	c := &Config{}
	err = yaml.Unmarshal(configFile, c)
	if err != nil {
		log.Fatal("failed to read config")
	}
	return c
}
