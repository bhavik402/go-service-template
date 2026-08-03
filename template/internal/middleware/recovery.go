package middleware

import (
	"log/slog"
	"net/http"

	"{{ module_name }}/internal/apperror"
	"{{ module_name }}/internal/util"
)

// Recovery turns a panic anywhere downstream into a logged 500 response
// instead of crashing the process.
func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered",
						slog.Any("panic", rec),
						slog.String("path", r.URL.Path),
						slog.String("request_id", RequestIDFromContext(r.Context())),
					)
					util.WriteError(w, apperror.Internal("internal server error", nil))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
