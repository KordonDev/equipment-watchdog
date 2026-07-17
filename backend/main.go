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
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kordondev/equipment-watchdog/changes"
	"github.com/kordondev/equipment-watchdog/config"
	"github.com/kordondev/equipment-watchdog/equipment"
	"github.com/kordondev/equipment-watchdog/gloveids"
	"github.com/kordondev/equipment-watchdog/members"
	"github.com/kordondev/equipment-watchdog/orders"
	"github.com/kordondev/equipment-watchdog/registrationcodes"
	"github.com/kordondev/equipment-watchdog/security"
	"github.com/kordondev/equipment-watchdog/tasks"
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

	if err := runMigrations(db); err != nil {
		panic(fmt.Sprintf("migration failed: %v", err))
	}

	userDB := users.NewDatebase(db)

	jwtService := security.NewJwtService(configuration.Origin, configuration.JwtSecret)

	router := gin.Default()
	_ = router.SetTrustedProxies(nil)

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{configuration.Origin}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	api := router.Group("/api")

	api.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	authorizeJWTMiddleware := security.AuthorizeJWTMiddleware(configuration.Origin, jwtService, configuration.Domain)
	userService := users.NewUserService(userDB, jwtService)
	webAuthNService, err := security.NewWebAuthNService(userService, configuration.Origin, configuration.Domain, jwtService, db)
	if err != nil {
		panic(fmt.Sprintf("Error creating webAuthn: %v", err))
	}
	_ = security.NewController(api, webAuthNService, configuration.Domain, authorizeJWTMiddleware)
	api.Use(authorizeJWTMiddleware)

	changeWriter := changes.NewChangeWriterService(db, userService)

	gloveIdService := gloveids.NewGloveIdService(db)

	equipmentService := equipment.NewEquipmentService(db, gloveIdService)
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

	gloveids.NewController(api, gloveIdService)

	taskService := tasks.NewTaskService(db)
	tasks.NewController(api, taskService)

	_ = router.Run(":8080")
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

func runMigrations(gormDB *gorm.DB) error {
	migrationPath := "file://migrations"
	if _, err := os.Stat("migrations"); os.IsNotExist(err) {
		migrationPath = "file:///migrations"
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("get underlying db: %w", err)
	}

	driver, err := sqlite3.WithInstance(sqlDB, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("create migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationPath, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
