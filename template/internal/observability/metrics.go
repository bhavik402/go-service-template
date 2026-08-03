package observability

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds the process-wide Prometheus collectors for HTTP traffic.
// Repositories, clients, {% if use_kafka %}messaging, {% endif %}{% if use_redis %}cache, {% endif %}etc. can register their own
// collectors the same way (promauto.With(reg).New...) if they need
// domain-specific signals later.
type Metrics struct {
	registry        *prometheus.Registry
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewMetrics() *Metrics {
	reg := prometheus.NewRegistry()
	reg.MustRegister(prometheus.NewGoCollector(), prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))

	return &Metrics{
		registry: reg,
		requestsTotal: promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests processed, labeled by route, method, and status.",
		}, []string{"route", "method", "status"}),
		requestDuration: promauto.With(reg).NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds.",
			Buckets: prometheus.DefBuckets,
		}, []string{"route", "method"}),
	}
}

// Handler exposes the /metrics scrape endpoint.
func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

// Observe records one completed HTTP request. Call it from middleware.
func (m *Metrics) Observe(route, method string, status int, duration time.Duration) {
	m.requestsTotal.WithLabelValues(route, method, strconv.Itoa(status)).Inc()
	m.requestDuration.WithLabelValues(route, method).Observe(duration.Seconds())
}
