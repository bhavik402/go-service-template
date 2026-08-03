package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"{{ module_name }}/internal/middleware"
	"{{ module_name }}/internal/observability"
)

type RouterConfig struct {
	UserHandler   *UserHandler
	HealthChecker *observability.HealthChecker
	Metrics       *observability.Metrics
	Logger        *slog.Logger
}

// NewRouter wires every entry point this service exposes today (HTTP
// only) onto the same handlers that {% if use_kafka %}a Kafka consumer or {% endif %}a cron job would call — the
// router is just one adapter in front of service/.
func NewRouter(cfg RouterConfig) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recovery(cfg.Logger))
	r.Use(middleware.Logging(cfg.Logger, cfg.Metrics))
	r.Use(middleware.CORS)

	r.Get("/health", cfg.HealthChecker.Live)
	r.Get("/ready", cfg.HealthChecker.Ready)
	r.Handle("/metrics", cfg.Metrics.Handler())

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", cfg.UserHandler.Create)
			r.Get("/", cfg.UserHandler.List)
			r.Get("/{id}", cfg.UserHandler.Get)
			r.Put("/{id}", cfg.UserHandler.Update)
			r.Delete("/{id}", cfg.UserHandler.Delete)
		})
	})

	return r
}
