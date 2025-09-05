package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MilliTime time.Time

func (t MilliTime) MarshalJSON() ([]byte, error) {
	RFC3339Milli := "2006-01-02T15:04:05.999Z07:00"
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(RFC3339Milli))
	return []byte(stamp), nil
}

type Response struct {
	Success     bool        `json:"success"`
	Data        interface{} `json:"data,omitempty"`
	Message     string      `json:"message,omitempty"`
	CurrentTime MilliTime   `json:"curTime"`
	StatusCode  int         `json:"-"`
}

func newResponse(data interface{}, success bool, message string, statusCode int) Response {
	return Response{
		Success:     success,
		Data:        data,
		Message:     message,
		CurrentTime: MilliTime(time.Now()),
		StatusCode:  statusCode,
	}
}

// WithPayload sends successful response with data
func WithPayload(w http.ResponseWriter, r *http.Request, data interface{}) {
	resp := newResponse(data, true, "", http.StatusOK)
	writeResponse(w, resp)
}

// WithMessage sends successful response with message only
func WithMessage(w http.ResponseWriter, r *http.Request, message string) {
	resp := newResponse(nil, true, message, http.StatusOK)
	writeResponse(w, resp)
}

// WithError sends error response with message
func WithError(w http.ResponseWriter, r *http.Request, message string) {
	resp := newResponse(nil, false, message, http.StatusBadRequest)
	writeResponse(w, resp)
}

// WithErrorPayload sends error response with data and message
func WithErrorPayload(w http.ResponseWriter, r *http.Request, message string, data interface{}) {
	resp := newResponse(data, false, message, http.StatusBadRequest)
	writeResponse(w, resp)
}

// WithStatusCode sends response with custom status code
func WithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	resp := newResponse(nil, false, message, statusCode)
	writeResponse(w, resp)
}

// WithUnauthorized sends unauthorized error response
func WithUnauthorized(w http.ResponseWriter, r *http.Request, message string) {
	if message == "" {
		message = "Unauthorized access"
	}
	resp := newResponse(nil, false, message, http.StatusUnauthorized)
	writeResponse(w, resp)
}

// WithForbidden sends forbidden error response
func WithForbidden(w http.ResponseWriter, r *http.Request, message string) {
	if message == "" {
		message = "Insufficient permissions"
	}
	resp := newResponse(nil, false, message, http.StatusForbidden)
	writeResponse(w, resp)
}

// WithInternalError sends internal server error response
func WithInternalError(w http.ResponseWriter, r *http.Request, message string) {
	if message == "" {
		message = "Internal server error"
	}
	resp := newResponse(nil, false, message, http.StatusInternalServerError)
	writeResponse(w, resp)
}

// WithNotFound sends not found error response
func WithNotFound(w http.ResponseWriter, r *http.Request, message string) {
	if message == "" {
		message = "Resource not found"
	}
	resp := newResponse(nil, false, message, http.StatusNotFound)
	writeResponse(w, resp)
}

// WithValidationError sends validation error response
func WithValidationError(w http.ResponseWriter, r *http.Request, message string) {
	if message == "" {
		message = "Validation failed"
	}
	resp := newResponse(nil, false, message, http.StatusUnprocessableEntity)
	writeResponse(w, resp)
}

// WithConflict sends conflict error response
func WithConflict(w http.ResponseWriter, r *http.Request, message string) {
	if message == "" {
		message = "Resource conflict"
	}
	resp := newResponse(nil, false, message, http.StatusConflict)
	writeResponse(w, resp)
}

// Error sends error response with custom status code (for backward compatibility)
func Error(w http.ResponseWriter, statusCode int, message string) {
	resp := newResponse(nil, false, message, statusCode)
	writeResponse(w, resp)
}

// JSON sends successful response with data and custom status code (for backward compatibility)
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	resp := newResponse(data, true, "", statusCode)
	writeResponse(w, resp)
}

// writeResponse writes the response to the client
func writeResponse(w http.ResponseWriter, resp Response) {
	w.Header().Set("Content-Type", "application/json")

	if resp.StatusCode != 0 {
		w.WriteHeader(resp.StatusCode)
	}

	json.NewEncoder(w).Encode(resp)
}
