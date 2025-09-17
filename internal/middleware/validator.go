package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validator.Struct(s)
}

func (v *Validator) ValidateJSON(next http.HandlerFunc, target interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(target); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		if err := v.ValidateStruct(target); err != nil {
			validationErrors := v.FormatValidationErrors(err)
			response := map[string]interface{}{
				"error":   "Validation failed",
				"details": validationErrors,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(response)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (v *Validator) FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				errors[field] = "This field is required"
			case "email":
				errors[field] = "Must be a valid email address"
			case "min":
				errors[field] = "Must be at least " + e.Param() + " characters"
			case "max":
				errors[field] = "Must be no more than " + e.Param() + " characters"
			default:
				errors[field] = "Invalid value"
			}
		}
	}

	return errors
}
