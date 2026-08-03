package observability

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
)

// Pinger is implemented by anything /ready should check before reporting
// the service healthy (a database pool, a cache client, ...).
type Pinger interface {
	Ping(ctx context.Context) error
}

// HealthChecker answers /health (liveness: the process is up) and /ready
// (readiness: the process can serve traffic) requests.
type HealthChecker struct {
	mu       sync.RWMutex
	checkers map[string]Pinger
}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{checkers: make(map[string]Pinger)}
}

// Register adds a named dependency that must respond for /ready to pass.
func (h *HealthChecker) Register(name string, p Pinger) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checkers[name] = p
}

// Live always returns 200 once the process is running.
func (h *HealthChecker) Live(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Ready pings every registered dependency and reports 200 only if all of
// them succeed.
func (h *HealthChecker) Ready(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	results := make(map[string]string, len(h.checkers))
	allOK := true

	for name, p := range h.checkers {
		if err := p.Ping(r.Context()); err != nil {
			results[name] = err.Error()
			allOK = false
			continue
		}
		results[name] = "ok"
	}

	status := http.StatusOK
	if !allOK {
		status = http.StatusServiceUnavailable
	}
	writeJSON(w, status, map[string]any{"status": statusLabel(allOK), "checks": results})
}

func statusLabel(ok bool) string {
	if ok {
		return "ok"
	}
	return "unavailable"
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
