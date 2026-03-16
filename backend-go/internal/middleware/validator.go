package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidationError represents a field validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateAndDecode reads JSON body, decodes into dest, and validates.
// Returns the decoded struct or writes a 400 error and returns nil.
func ValidateAndDecode(w http.ResponseWriter, r *http.Request, dest interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":   "Datos de entrada inválidos",
			"details": []ValidationError{{Field: "body", Message: fmt.Sprintf("JSON inválido: %v", err)}},
		})
		return false
	}

	if err := validate.Struct(dest); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			details := make([]ValidationError, len(validationErrors))
			for i, ve := range validationErrors {
				details[i] = ValidationError{
					Field:   ve.Field(),
					Message: fmt.Sprintf("failed on '%s' validation", ve.Tag()),
				}
			}
			WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
				"error":   "Datos de entrada inválidos",
				"details": details,
			})
			return false
		}
		WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": "Datos de entrada inválidos",
		})
		return false
	}

	return true
}
