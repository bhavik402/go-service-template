package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"{{ module_name }}/internal/observability"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// Logging writes one structured access log line per request and, if
// metrics is non-nil, records the request in Prometheus. The route label
// uses the matched chi pattern (e.g. /users/{id}) rather than the raw
// path, so metrics don't explode in cardinality with every distinct ID.
func Logging(logger *slog.Logger, metrics *observability.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(rec, r)

			duration := time.Since(start)
			route := routePattern(r)

			logger.Info("request handled",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("route", route),
				slog.Int("status", rec.status),
				slog.Duration("duration", duration),
				slog.String("request_id", RequestIDFromContext(r.Context())),
			)

			if metrics != nil {
				metrics.Observe(route, r.Method, rec.status, duration)
			}
		})
	}
}

func routePattern(r *http.Request) string {
	if rctx := chi.RouteContext(r.Context()); rctx != nil {
		if pattern := rctx.RoutePattern(); pattern != "" {
			return pattern
		}
	}
	return r.URL.Path
}
