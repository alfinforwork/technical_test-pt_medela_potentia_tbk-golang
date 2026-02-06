package main

import (
	"strconv"
	"technical-test/src/config"
	"technical-test/src/database"
	"technical-test/src/handler"
	"technical-test/src/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"gorm.io/gorm"
)

func main() {
	app := setupFiberApp()
	app.Use(logger.New())
	db := setupDatabase()
	defer closeDatabase(db)

	// Setup routes
	handler.Routes(app, db)

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
