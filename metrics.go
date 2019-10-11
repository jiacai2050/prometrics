package prometrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// runMetricsCollector / fallbackMetricsCollector share those metric
	successDuration                *prometheus.HistogramVec
	errFailureDuration             *prometheus.HistogramVec
	errConcurrencyLimitRejectTotal *prometheus.CounterVec

	// runMetricsCollector's metrics
	errTimeoutDuration    *prometheus.HistogramVec
	errBadRequestDuration *prometheus.HistogramVec
	errShortCircuitTotal  *prometheus.CounterVec

	// circuitMetricsCollector's metrics
	closedTotal *prometheus.CounterVec
	openedTotal *prometheus.CounterVec
)

const (
	ns               = "circuit"
	circuitNameLabel = "name"
	funcTypeLabel    = "func"

	runCmd      = "run"
	fallbackCmd = "fallback"
)

func init() {
	successDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: ns,
			Name:      "success_duration_seconds",
			Help:      "Total of successful func run",
		},
		[]string{funcTypeLabel, circuitNameLabel},
	)

	errFailureDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: ns,
			Name:      "failure_duration_seconds",
			Help:      "Duration of failed func run",
		},
		[]string{funcTypeLabel, circuitNameLabel},
	)

	errConcurrencyLimitRejectTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: ns,
			Name:      "concurrency_reject_total",
			Help:      "Total reject requests for reach concurrency limit",
		},
		[]string{funcTypeLabel, circuitNameLabel},
	)

	errTimeoutDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: ns,
			Name:      "timeout_duration_seconds",
			Help:      "Duration of timeout func run",
		},
		[]string{circuitNameLabel},
	)

	errBadRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: ns,
			Name:      "bad_request_duration_seconds",
			Help:      "Duration of bad request request",
		},
		[]string{circuitNameLabel},
	)

	errShortCircuitTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: ns,
			Name:      "short_circuit_total",
			Help:      "Total of runFunc is not called because the circuit was open",
		},
		[]string{circuitNameLabel},
	)

	closedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: ns,
			Name:      "closed_total",
			Help:      "Total of circuit transitions from open to closed",
		},
		[]string{circuitNameLabel},
	)

	openedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: ns,
			Name:      "opened_total",
			Help:      "Total of circuit transitions from closed to opened",
		},
		[]string{circuitNameLabel},
	)
}
