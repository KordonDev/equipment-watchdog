package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/security"
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
	args := parseArguments()

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

	router.LoadHTMLGlob("templates/*.html")

	router.GET("/index/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"name": name,
		})
	})

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

type CmdArgs struct {
	Debug  bool
	Domain string
	Origin string
}

func parseArguments() *CmdArgs {
	debug := flag.Bool("debug", false, "log debug information")
	domain := flag.String("domain", "localhost", "Generally the domain name for your site with webAuthn")
	origin := flag.String("origin", "http://localhost:8080", "The origin URL for WebAuthn requests")
	flag.Parse()
	fmt.Printf("domain %s, origin %s \n", *domain, *origin)
	return &CmdArgs{
		Debug:  *debug,
		Domain: *domain,
		Origin: *origin,
	}
}
