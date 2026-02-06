package handler

import (
	"technical-test/src/response"
	"technical-test/src/usecase"
	"technical-test/src/utils"

	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with name, email, and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body object{name=string,email=string,password=string} true "Register Request"
// @Success 200 {object} response.ResponseSuccess "User registered successfully"
// @Failure 400 {object} response.ResponseError "Validation error" example("message":"Invalid input")
// @Router /v1/auth/register [post]
func (h *AuthHandler) Register(c fiber.Ctx) error {
	var body struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	if err := c.Bind().Body(&body); err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, utils.FormatValidationError(err), nil)
	}

	user, err := h.authUsecase.Register(body.Name, body.Email, body.Password)
	if err != nil {
		return response.Error(c, err.Error(), nil)
	}

	return response.Success(c, "User registered successfully", fiber.Map{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	}, nil)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with email and password, returns JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body object{email=string,password=string} true "Login Request"
// @Success 200 {object} response.ResponseSuccess{data=response.LoginResponse} "Login successful"
// @Failure 400 {object} response.ResponseError "Validation error"
// @Failure 401 {object} response.ResponseError "Invalid credentials"
// @Router /v1/auth/login [post]
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var body struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.Bind().Body(&body); err != nil {
		c.Status(fiber.StatusBadRequest)
		return response.Error(c, utils.FormatValidationError(err), nil)
	}

	token, user, err := h.authUsecase.Login(body.Email, body.Password)
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
