package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) Validate(data interface{}) error {
	return v.validate.Struct(data)
}

func FormatValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				messages = append(messages, fmt.Sprintf("%s is required", field))
			case "email":
				messages = append(messages, fmt.Sprintf("%s must be a valid email", field))
			case "min":
				messages = append(messages, fmt.Sprintf("%s must be at least %s characters", field, e.Param()))
			case "max":
				messages = append(messages, fmt.Sprintf("%s must be at most %s characters", field, e.Param()))
			default:
				messages = append(messages, fmt.Sprintf("%s is invalid", field))
			}
		}
		return strings.Join(messages, ", ")
	}
	return err.Error()
}
