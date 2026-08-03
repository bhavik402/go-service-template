// Package util holds small helpers shared across handlers: writing JSON
// responses, mapping apperror codes to HTTP statuses, and pagination
// parsing. Keep it small — anything that grows its own concept deserves
// its own package instead of living here.
package util

import (
	"encoding/json"
	"errors"
	"net/http"

	"{{ module_name }}/internal/apperror"
)

// ErrorResponse is the standard error envelope returned by every handler.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// WriteJSON writes body as a JSON response with the given status code.
func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}

// WriteError maps err to an HTTP status and writes the standard error
// envelope. Unrecognized errors are treated as internal errors and their
// details are never leaked to the client — log them at the call site.
func WriteError(w http.ResponseWriter, err error) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		WriteJSON(w, statusForCode(appErr.Code), ErrorResponse{
			Code:    string(appErr.Code),
			Message: appErr.Message,
		})
		return
	}

	WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
		Code:    string(apperror.CodeInternal),
		Message: "internal server error",
	})
}

func statusForCode(code apperror.Code) int {
	switch code {
	case apperror.CodeNotFound:
		return http.StatusNotFound
	case apperror.CodeConflict:
		return http.StatusConflict
	case apperror.CodeInvalidInput:
		return http.StatusBadRequest
	case apperror.CodeUnauthorized:
		return http.StatusUnauthorized
	case apperror.CodeForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
