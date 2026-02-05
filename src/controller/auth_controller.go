package controller

import (
	"technical-test/src/response"
	"technical-test/src/service"
	"technical-test/src/utils"

	"github.com/gofiber/fiber/v3"
)

type AuthController struct {
	AuthService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

func (ac *AuthController) Register(c fiber.Ctx) error {
	var body struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	if err := c.Bind().Body(&body); err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, utils.FormatValidationError(err), nil)
	}

	user, err := ac.AuthService.Register(body.Name, body.Email, body.Password)
	if err != nil {
		return response.Error(c, err.Error(), nil)
	}

	return response.Success(c, "User registered successfully", fiber.Map{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	}, nil)
}

func (ac *AuthController) Login(c fiber.Ctx) error {
	var body struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.Bind().Body(&body); err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, utils.FormatValidationError(err), nil)
	}

	token, user, err := ac.AuthService.Login(body.Email, body.Password)
	if err != nil {
		return response.Error(c, err.Error(), nil)
	}

	return response.Success(c, "Login successful", fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	}, nil)
}
