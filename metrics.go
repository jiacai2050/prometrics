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
	opened      *prometheus.GaugeVec
)

const (
	subSystem        = "circuit"
	circuitNameLabel = "name"
	funcTypeLabel    = "func"

	runCmd      = "run"
	fallbackCmd = "fallback"
)

func initMetrics(namespace string) {
	successDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subSystem,
			Name:      "success_duration_seconds",
			Help:      "Duration of successful func run",
		},
		[]string{funcTypeLabel, circuitNameLabel},
	)

	errFailureDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subSystem,
			Name:      "failure_duration_seconds",
			Help:      "Duration of failed func run",
		},
		[]string{funcTypeLabel, circuitNameLabel},
	)

	errConcurrencyLimitRejectTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subSystem,
			Name:      "concurrency_reject_total",
			Help:      "Total reject requests for reach concurrency limit",
		},
		[]string{funcTypeLabel, circuitNameLabel},
	)

	errTimeoutDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subSystem,
			Name:      "timeout_duration_seconds",
			Help:      "Duration of timeout func run",
		},
		[]string{circuitNameLabel},
	)

	errBadRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subSystem,
			Name:      "bad_request_duration_seconds",
			Help:      "Duration of bad request request",
		},
		[]string{circuitNameLabel},
	)

	errShortCircuitTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subSystem,
			Name:      "short_circuit_total",
			Help:      "Total of runFunc is not called because the circuit was open",
		},
		[]string{circuitNameLabel},
	)

	closedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subSystem,
			Name:      "closed_total",
			Help:      "Total of circuit transitions from open to closed",
		},
		[]string{circuitNameLabel},
	)

	openedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subSystem,
			Name:      "opened_total",
			Help:      "Total of circuit transitions from closed to opened",
		},
		[]string{circuitNameLabel},
	)

	opened = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subSystem,
			Name:      "opened",
			Help:      "The status of a circuit, 1 opened, 0 closed",
		},
		[]string{circuitNameLabel},
	)
}
