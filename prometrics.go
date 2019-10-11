package prometrics

import (
	"sync"
	"time"

	"github.com/cep21/circuit/v3"
	"github.com/prometheus/client_golang/prometheus"
)

type CommandFactory struct {
}

type runMetricsCollector struct {
	name string
}

// Success sends a success to prometheus
func (c *runMetricsCollector) Success(now time.Time, duration time.Duration) {
	successDuration.WithLabelValues(runCmd, c.name).Observe(duration.Seconds())
}

// ErrFailure sends a failure to prometheus
func (c *runMetricsCollector) ErrFailure(now time.Time, duration time.Duration) {
	errFailureDuration.WithLabelValues(runCmd, c.name).Observe(duration.Seconds())
}

// ErrTimeout sends a timeout to prometheus
func (c *runMetricsCollector) ErrTimeout(now time.Time, duration time.Duration) {
	errTimeoutDuration.WithLabelValues(c.name).Observe(duration.Seconds())
}

// ErrBadRequest sends a bad request error to prometheus
func (c *runMetricsCollector) ErrBadRequest(now time.Time, duration time.Duration) {
	errBadRequestDuration.WithLabelValues(c.name).Observe(duration.Seconds())
}

// ErrInterrupt sends an interrupt error to prometheus
func (c *runMetricsCollector) ErrInterrupt(now time.Time, duration time.Duration) {
	// no need to implement
}

// ErrShortCircuit sends a short circuit to prometheus
func (c *runMetricsCollector) ErrShortCircuit(now time.Time) {
	errShortCircuitTotal.WithLabelValues(c.name).Inc()
}

// ErrConcurrencyLimitReject sends a concurrency limit error to prometheus
func (c *runMetricsCollector) ErrConcurrencyLimitReject(now time.Time) {
	errConcurrencyLimitRejectTotal.WithLabelValues(runCmd, c.name).Inc()
}

var _ circuit.RunMetrics = (*runMetricsCollector)(nil)

type fallbackMetricsCollector struct {
	name string
}

// Success sends a success to prometheus
func (c *fallbackMetricsCollector) Success(now time.Time, duration time.Duration) {
	successDuration.WithLabelValues(fallbackCmd, c.name).Observe(duration.Seconds())
}

// ErrConcurrencyLimitReject sends a concurrency-limit to prometheus
func (c *fallbackMetricsCollector) ErrConcurrencyLimitReject(now time.Time) {
	errConcurrencyLimitRejectTotal.WithLabelValues(fallbackCmd, c.name).Inc()
}

// ErrFailure sends a failure to prometheus
func (c *fallbackMetricsCollector) ErrFailure(now time.Time, duration time.Duration) {
	errFailureDuration.WithLabelValues(fallbackCmd, c.name).Observe(duration.Seconds())
}

var _ circuit.FallbackMetrics = (*fallbackMetricsCollector)(nil)

type circuitMetricsCollector struct {
	name string
}

// Closed sets a gauge as closed for the collector
func (c *circuitMetricsCollector) Closed(now time.Time) {
	closedTotal.WithLabelValues(c.name).Inc()
}

// Opened sets a gauge as opened for the collector
func (c *circuitMetricsCollector) Opened(now time.Time) {
	openedTotal.WithLabelValues(c.name).Inc()
}

var _ circuit.Metrics = (*circuitMetricsCollector)(nil)

func (c *CommandFactory) CommandProperties(circuitName string) circuit.Config {
	return circuit.Config{
		Metrics: circuit.MetricsCollectors{
			Run:      []circuit.RunMetrics{&runMetricsCollector{circuitName}},
			Fallback: []circuit.FallbackMetrics{&fallbackMetricsCollector{circuitName}},
			Circuit:  []circuit.Metrics{&circuitMetricsCollector{circuitName}},
		},
	}
}

var (
	once sync.Once
	inst *CommandFactory
)

func GetFactory(r prometheus.Registerer) *CommandFactory {
	once.Do(func() {
		r.MustRegister(
			successDuration,
			errFailureDuration,
			errConcurrencyLimitRejectTotal,
			errTimeoutDuration,
			errBadRequestDuration,
			errShortCircuitTotal,
			closedTotal,
			openedTotal,
		)
		inst = new(CommandFactory)
	})

	return inst
}
