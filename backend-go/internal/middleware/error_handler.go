package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

// AppError is a structured HTTP error.
type AppError struct {
	StatusCode int         `json:"-"`
	Message    string      `json:"error"`
	Details    interface{} `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new AppError.
func NewAppError(statusCode int, message string) *AppError {
	return &AppError{StatusCode: statusCode, Message: message}
}

// WriteError writes a JSON error response.
func WriteError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*AppError); ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appErr.StatusCode)
		json.NewEncoder(w).Encode(appErr)
		return
	}

	log.Printf("Unhandled error: %v", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "Error interno del servidor",
	})
}

// WriteJSON writes a JSON response with the given status code.
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// ReadJSON reads and decodes JSON from request body into dest.
func ReadJSON(r *http.Request, dest interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(dest)
}
