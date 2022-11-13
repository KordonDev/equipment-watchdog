package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/security"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type member struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

var members = []member{
	{Id: "1", Name: "Arne"},
	{Id: "2", Name: "Luka"},
}

func getMembers(c *gin.Context) {
	c.JSON(http.StatusOK, members)
}

func addMember(c *gin.Context) {
	var newMember member

	if err := c.BindJSON(&newMember); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	fmt.Println(newMember)

	members = append(members, newMember)
	c.JSON(http.StatusCreated, newMember)
}

func main() {
	args := parseConfig()

	router := gin.Default()
	api := router.Group("/api")

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:5173"}
	corsConfig.AllowCredentials = true
	api.Use(cors.New(corsConfig))

	members := api.Group("/members", security.AuthorizeJWTMiddleware())
	members.GET("/", getMembers)
	members.POST("/", addMember)

	db := createDB(args.Debug)
	userDB := security.NewUserDB(db)
	webAuthNService := security.NewWebAuthNService(userDB, args.Origin, args.Domain)

	api.GET("/register/:username", webAuthNService.StartRegister)
	api.POST("/register/:username", webAuthNService.FinishRegistration)

	api.GET("/login/:username", webAuthNService.StartLogin)
	api.POST("/login/:username", webAuthNService.FinishLogin)
	api.POST("/logout", webAuthNService.Logout)

	router.Run("localhost:8080")
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
	Debug  bool
	Domain string
	Origin string
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
