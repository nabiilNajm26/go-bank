package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ValidationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}

func ValidateRequest(payload interface{}) error {
	return validate.Struct(payload)
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	if validationErr, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)
		for _, err := range validationErr {
			errors[err.Field()] = getErrorMessage(err)
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Validation failed",
			"fields": errors,
		})
	}
	return err
}

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	case "uuid":
		return "Invalid UUID format"
	case "gt":
		return "Value must be greater than " + err.Param()
	default:
		return "Invalid value"
	}
}