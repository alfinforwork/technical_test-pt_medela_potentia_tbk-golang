package router

import (
	"technical-test/src/controller"
	"technical-test/src/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetupAuthRoutes(router fiber.Router, db *gorm.DB, as *service.AuthService) {
	authGroup := router.Group("/auth")
	authController := controller.NewAuthController(*as)

	authGroup.Post("/register", authController.Register)
	authGroup.Post("/login", authController.Login)
}
