package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type ReplicationMetrics struct {
	// HTTP
	HTTPRequestsTotal    *prometheus.CounterVec
	HTTPRequestDuration  *prometheus.HistogramVec
	HTTPInflightRequests *prometheus.GaugeVec

	// Replication phase latency
	MasterWriteDuration   *prometheus.HistogramVec
	FollowerWriteDuration *prometheus.HistogramVec
	SyncOverheadDuration  *prometheus.HistogramVec

	// Replication outcomes
	FollowerWriteTotal    *prometheus.CounterVec
	PartialCommitTotal    *prometheus.CounterVec
	FollowerTimeoutsTotal *prometheus.CounterVec
}

func NewReplicationMetrics(reg prometheus.Registerer) *ReplicationMetrics {
	m := &ReplicationMetrics{
		HTTPRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total HTTP requests processed",
			},
			[]string{"route", "method", "status"},
		),

		HTTPRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"route", "method"},
		),

		HTTPInflightRequests: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "http_inflight_requests",
				Help: "Number of in-flight HTTP requests",
			},
			[]string{"route"},
		),

		MasterWriteDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "replication_master_write_duration_seconds",
				Help:    "Duration of master write in seconds",
				Buckets: []float64{0.001, 0.003, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2, 5},
			},
			[]string{"op"},
		),

		FollowerWriteDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "replication_follower_write_duration_seconds",
				Help:    "Duration of follower write in seconds",
				Buckets: []float64{0.001, 0.003, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2, 5},
			},
			[]string{"op", "follower"},
		),

		SyncOverheadDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "replication_sync_overhead_duration_seconds",
				Help:    "Overhead of synchronous replication (time spent waiting for followers etc.)",
				Buckets: []float64{0.001, 0.003, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2, 5},
			},
			[]string{"op"},
		),

		FollowerWriteTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "replication_follower_write_total",
				Help: "Follower write outcomes",
			},
			[]string{"op", "follower", "result"}, // ok|error|timeout
		),

		PartialCommitTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "replication_partial_commit_total",
				Help: "Master committed but follower failed (sync request must fail).",
			},
			[]string{"op"},
		),

		FollowerTimeoutsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "replication_follower_timeouts_total",
				Help: "Follower timeouts",
			},
			[]string{"op", "follower"},
		),
	}

	reg.MustRegister(
		m.HTTPRequestsTotal,
		m.HTTPRequestDuration,
		m.HTTPInflightRequests,
		m.MasterWriteDuration,
		m.FollowerWriteDuration,
		m.SyncOverheadDuration,
		m.FollowerWriteTotal,
		m.PartialCommitTotal,
		m.FollowerTimeoutsTotal,
	)

	return m
}

// ---- HTTP middleware ----

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (m *ReplicationMetrics) Instrument(route string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}

		m.HTTPInflightRequests.WithLabelValues(route).Inc()
		start := time.Now()

		defer func() {
			m.HTTPInflightRequests.WithLabelValues(route).Dec()

			dur := time.Since(start).Seconds()
			m.HTTPRequestDuration.WithLabelValues(route, r.Method).Observe(dur)
			m.HTTPRequestsTotal.WithLabelValues(route, r.Method, strconv.Itoa(sw.status)).Inc()
		}()

		next.ServeHTTP(sw, r)
	})
}
