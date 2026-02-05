package config

import (
	"fmt"

	"github.com/gofiber/fiber/v3/log"
	"github.com/spf13/viper"
)

var (
	IsProd              bool
	AppHost             string
	AppPort             int
	DBHost              string
	DBUser              string
	DBPassword          string
	DBName              string
	DBPort              int
	DBMigrate           bool
	JWTSecret           string
	JWTAccessExp        int
	JWTRefreshExp       int
	JWTResetPasswordExp int
	JWTVerifyEmailExp   int
)

func init() {
	setDefaults()
	loadConfig()

	// server configuration
	IsProd = viper.GetString("APP_ENV") == "prod"
	AppPort = viper.GetInt("APP_PORT")

	// database configuration
	DBHost = viper.GetString("DB_HOST")
	DBUser = viper.GetString("DB_USER")
	DBPassword = viper.GetString("DB_PASSWORD")
	DBName = viper.GetString("DB_NAME")
	DBPort = viper.GetInt("DB_PORT")
	DBMigrate = viper.GetBool("DB_MIGRATE")

	// jwt configuration
	JWTSecret = viper.GetString("JWT_SECRET")
	JWTAccessExp = viper.GetInt("JWT_ACCESS_EXP_MINUTES")
	JWTRefreshExp = viper.GetInt("JWT_REFRESH_EXP_DAYS")
	JWTResetPasswordExp = viper.GetInt("JWT_RESET_PASSWORD_EXP_MINUTES")
	JWTVerifyEmailExp = viper.GetInt("JWT_VERIFY_EMAIL_EXP_MINUTES")

	fmt.Printf("PORT: %d \n", AppPort)
}

func setDefaults() {
	// this function to support docker env variable which not use .env file
	// Set default values
	viper.SetDefault("APP_ENV", "dev")
	viper.SetDefault("APP_PORT", 3000)
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_USER", "root")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_NAME", "test")
	viper.SetDefault("DB_PORT", 3306)
	viper.SetDefault("DB_MIGRATE", false)
	viper.SetDefault("JWT_SECRET", "your-secret-key")
	viper.SetDefault("JWT_ACCESS_EXP_MINUTES", 15)
	viper.SetDefault("JWT_REFRESH_EXP_DAYS", 7)
	viper.SetDefault("JWT_RESET_PASSWORD_EXP_MINUTES", 30)
	viper.SetDefault("JWT_VERIFY_EMAIL_EXP_MINUTES", 60)

	// Bind environment variables
	viper.BindEnv("APP_ENV", "APP_ENV")
	viper.BindEnv("APP_PORT", "APP_PORT")
	viper.BindEnv("DB_HOST", "DB_HOST")
	viper.BindEnv("DB_USER", "DB_USER")
	viper.BindEnv("DB_PASSWORD", "DB_PASSWORD")
	viper.BindEnv("DB_NAME", "DB_NAME")
	viper.BindEnv("DB_PORT", "DB_PORT")
	viper.BindEnv("DB_MIGRATE", "DB_MIGRATE")
	viper.BindEnv("JWT_SECRET", "JWT_SECRET")
	viper.BindEnv("JWT_ACCESS_EXP_MINUTES", "JWT_ACCESS_EXP_MINUTES")
	viper.BindEnv("JWT_REFRESH_EXP_DAYS", "JWT_REFRESH_EXP_DAYS")
	viper.BindEnv("JWT_RESET_PASSWORD_EXP_MINUTES", "JWT_RESET_PASSWORD_EXP_MINUTES")
	viper.BindEnv("JWT_VERIFY_EMAIL_EXP_MINUTES", "JWT_VERIFY_EMAIL_EXP_MINUTES")
}

func loadConfig() {
	configPaths := []string{
		"./",     // For app
		"../../", // For test folder
	}

	for _, path := range configPaths {
		viper.SetConfigFile(path + ".env")

		if err := viper.ReadInConfig(); err == nil {
			log.Infof("Config file loaded from %s.env", path)
			return
		}
	}

	log.Warn("No .env file found, using default values and environment variables")
}
