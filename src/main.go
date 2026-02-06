package main

import (
	"strconv"
	"technical-test/docs"
	"technical-test/src/config"
	"technical-test/src/database"
	"technical-test/src/middleware"
	"technical-test/src/routes"
	"technical-test/src/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"gorm.io/gorm"
)

// Swagger documentation initialization
var _ = docs.SwaggerInfo

// @title Workflow Management API
// @version 1.0
// @description API for managing workflows, steps, and requests
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8000
// @basePath /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	app := setupFiberApp()
	app.Use(logger.New())
	db := setupDatabase()
	defer closeDatabase(db)

	// Setup Swagger routes
	app.Get("/swagger", middleware.SwaggerHandler())
	app.Get("/swagger.json", middleware.SwaggerHandler())

	// Setup routes
	routes.SetupRoutes(app, db)

	PORT := strconv.Itoa(config.AppPort)
	if PORT == "" {
		PORT = "3000"
	}
	if err := app.Listen(":" + PORT); err != nil {
		panic(err)
	}
}

func setupFiberApp() *fiber.App {
	app := fiber.New(fiber.Config{
		StructValidator: utils.NewValidator(),
	})
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	return app
}

func setupDatabase() *gorm.DB {
	db := database.Connect(config.DBHost, config.DBName)
	return db
}

func closeDatabase(db *gorm.DB) {
	sqlDB, errDB := db.DB()
	if errDB != nil {
		log.Errorf("Error getting database instance: %v", errDB)
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Errorf("Error closing database connection: %v", err)
	} else {
		log.Info("Database connection closed successfully")
	}
}
