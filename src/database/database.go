package database

import (
	"fmt"
	"technical-test/src/config"
	"technical-test/src/model"
	"time"

	"github.com/gofiber/fiber/v3/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(dbHost, dbName string) *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser, config.DBPassword, dbHost, config.DBPort, dbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		TranslateError:         true,
	})
	if err != nil {
		log.Errorf("Failed to connect to database: %+v", err)
		panic(fmt.Sprintf("Database connection failed: %v", err))
	}

	sqlDB, errDB := db.DB()
	if errDB != nil {
		log.Errorf("Failed to get database connection: %+v", errDB)
		panic(fmt.Sprintf("Failed to get database connection: %v", errDB))
	}

	// Config connection pooling
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	if config.DBMigrate {
		db.AutoMigrate(
			&model.User{},
			&model.Workflow{},
			&model.Step{},
			&model.Request{},
		)
	}
	return db
}
