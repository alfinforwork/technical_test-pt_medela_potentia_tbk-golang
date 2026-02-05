package response

import "github.com/gofiber/fiber/v3"

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
