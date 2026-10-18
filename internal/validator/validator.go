// Package validator via validator v10
package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator that wraps the validator from a third party library
type Validator struct {
	validate *validator.Validate
}

// New Validator contructor
func New() *Validator {
	return &Validator{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

// ValidationError struct for validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validate validate a struct by its tags and returns readable tags
func (v *Validator) Validate(s interface{}) []ValidationError {
	err := v.validate.Struct(s)

	if err == nil {
		return nil
	}

	var errors []ValidationError

	validationErrors, ok := err.(validator.ValidationErrors)

	if !ok {
		return nil
	}

	for _, e := range validationErrors {
		errors = append(errors, ValidationError{
			Field:   strings.ToLower(e.Field()),
			Message: translateTag(e),
		})
	}

	return errors
}

func translateTag(e validator.FieldError) string {
	var res string

	switch e.Tag() {
	case "required":
		res = "Field is required"
	case "email":
		res = "Invalid email"
	case "min":
		res = fmt.Sprintf("Must contains at least %s characters", e.Param())
	case "max":
		res = fmt.Sprintf("Must contains at most %s characters", e.Param())
	default:
		res = "Invalid value"
	}

	return res
}
