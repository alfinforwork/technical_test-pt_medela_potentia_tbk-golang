package middleware

import (
	"strings"
	"technical-test/src/response"
	"technical-test/src/service"

	"github.com/gofiber/fiber/v3"
)

func JWTProtected() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.Status(fiber.StatusUnauthorized)
			return response.Error(c, "Missing or invalid authorization header", nil)
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, claims, err := service.ParseToken(tokenString)
		if err != nil || token == nil || !token.Valid {
			c.Status(fiber.StatusUnauthorized)
			return response.Error(c, "Invalid token", nil)
		}

		if sub, ok := claims["sub"]; ok {
			c.Locals("user_id", sub)
		}
		if email, ok := claims["email"]; ok {
			c.Locals("email", email)
		}

		return c.Next()
	}
}
