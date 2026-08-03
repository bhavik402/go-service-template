// Package observability provides structured logging, metrics, and health
// checks — the signals the application itself emits. Operational concerns
// (dashboards, alert rules, scrape configs) live outside the codebase.
package observability

import (
	"log/slog"
	"os"
	"strings"
)

// NewLogger builds a JSON structured logger. level accepts debug, info,
// warn, or error (case-insensitive); anything else falls back to info.
func NewLogger(level string) *slog.Logger {
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn", "warning":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	})

	return slog.New(handler)
}
