package response

import "github.com/gofiber/fiber/v3"

// ========================================================
// Response Structures For Swagger Documentation
// ========================================================

// ResponseSuccess represents a successful API response
// @Description Success response object
type ResponseSuccess struct {
	Status  string      `json:"status" example:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// ResponseError represents an error API response
// @Description Error response object
type ResponseError struct {
	Status  string      `json:"status" example:"error"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// LoginResponse represents login response data
type LoginResponse struct {
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserResponse `json:"user"`
}

// UserResponse represents user data in response
type UserResponse struct {
	ID    uint   `json:"id" example:"1"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"user@example.com"`
}

// ========================================================
// Response Helpers
// ========================================================

func Success(c fiber.Ctx, message string, data interface{}, additional []interface{}) error {
	result := fiber.Map{
		"status":  "success",
		"data":    data,
		"message": message,
	}
	for i := 0; i < len(additional); i += 2 {
		if i+1 < len(additional) {
			result[additional[i].(string)] = additional[i+1]
		}
	}
	return c.JSON(result)
}

func Error(c fiber.Ctx, message string, additional []interface{}) error {
	code := c.Response().StatusCode()
	if code == fiber.StatusOK || code == 0 {
		code = fiber.StatusBadRequest
	}
	result := fiber.Map{
		"status":  "error",
		"data":    nil,
		"message": message,
	}
	for i := 0; i < len(additional); i += 2 {
		if i+1 < len(additional) {
			result[additional[i].(string)] = additional[i+1]
		}
	}
	return c.Status(code).JSON(result)
}
