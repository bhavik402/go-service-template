// Package middleware holds cross-cutting HTTP concerns: request IDs,
// panic recovery, access logging, CORS. Business logic never lives here —
// these wrap the handler chain, they don't replace it.
package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const requestIDKey contextKey = "request_id"

const RequestIDHeader = "X-Request-ID"

// RequestID assigns a request ID (reusing an inbound X-Request-ID header
// if present), stores it in the request context, and echoes it back on
// the response.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(RequestIDHeader)
		if id == "" {
			id = uuid.NewString()
		}

		w.Header().Set(RequestIDHeader, id)
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestIDFromContext returns the request ID stored by RequestID, or ""
// if none is present.
func RequestIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}
