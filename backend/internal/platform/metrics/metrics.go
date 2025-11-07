package metrics

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	HttpRequestsTotal     *prometheus.CounterVec
	HttpRequestDuration   *prometheus.HistogramVec
	HttpRequestsInFlight  prometheus.Gauge
	DbQueryDuration       *prometheus.HistogramVec
	DbQueriesTotal        *prometheus.CounterVec
	PermissionChecksTotal *prometheus.CounterVec
	ValidationErrorsTotal *prometheus.CounterVec
	DatabaseErrorsTotal   prometheus.Counter
}

var M *Metrics

func Init() {
	M = &Metrics{
		HttpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: "http_requests_total", Help: "Total HTTP requests"},
			[]string{"method", "endpoint", "status"},
		),
		HttpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{Name: "http_request_duration_seconds", Help: "HTTP request duration in seconds", Buckets: []float64{.001, .005, .01, .05, .1, .5, 1, 2, 5}},
			[]string{"method", "endpoint"},
		),
		HttpRequestsInFlight: prometheus.NewGauge(
			prometheus.GaugeOpts{Name: "http_requests_in_flight", Help: "Current number of HTTP requests in flight"},
		),
		DbQueryDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{Name: "db_query_duration_seconds", Help: "Database query duration in seconds", Buckets: []float64{.001, .01, .05, .1, .5, 1}},
			[]string{"operation", "table"},
		),
		DbQueriesTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: "db_queries_total", Help: "Total database queries"},
			[]string{"operation", "table", "status"},
		),
		PermissionChecksTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: "permission_checks_total", Help: "Total permission checks"},
			[]string{"resource", "action", "allowed"},
		),
		ValidationErrorsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: "validation_errors_total", Help: "Total validation errors"},
			[]string{"field", "error_type"},
		),
		DatabaseErrorsTotal: prometheus.NewCounter(
			prometheus.CounterOpts{Name: "database_errors_total", Help: "Total database errors"},
		),
	}

	prometheus.MustRegister(M.HttpRequestsTotal)
	prometheus.MustRegister(M.HttpRequestDuration)
	prometheus.MustRegister(M.HttpRequestsInFlight)
	prometheus.MustRegister(M.DbQueryDuration)
	prometheus.MustRegister(M.DbQueriesTotal)
	prometheus.MustRegister(M.PermissionChecksTotal)
	prometheus.MustRegister(M.ValidationErrorsTotal)
	prometheus.MustRegister(M.DatabaseErrorsTotal)

	log.Println("✓ Prometheus metrics initialized")
}
