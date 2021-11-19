package prometrics

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/cep21/circuit/v3"
	"github.com/cep21/circuit/v3/closers/hystrix"
	"github.com/cep21/circuit/v3/metrics/rolling"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

var (
	errRunFunc      = fmt.Errorf("err in run func")
	errFallbackFunc = fmt.Errorf("err in fallback func")
)

func init() {
	// curl 0:8080/metrics | grep circuit
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
}

func TestCommandFactory(t *testing.T) {
	f := rolling.StatFactory{}
	prom := GetFactory(prometheus.DefaultRegisterer)
	h := circuit.Manager{
		DefaultCircuitProperties: []circuit.CommandPropertiesConstructor{
			f.CreateConfig,
			prom.CommandProperties,
			func(circuitName string) circuit.Config {
				return circuit.Config{
					Execution: circuit.ExecutionConfig{
						MaxConcurrentRequests: 10,
					},
					General: circuit.GeneralConfig{
						OpenToClosedFactory: hystrix.CloserFactory(hystrix.ConfigureCloser{
							SleepWindow: time.Millisecond * 800,
						}),
						ClosedToOpenFactory: hystrix.OpenerFactory(hystrix.ConfigureOpener{
							RequestVolumeThreshold:   2,
							ErrorThresholdPercentage: 50,
						}),
					},
				}
			},
		},
	}

	cmdName := "test-cmd"
	c := h.MustCreateCircuit(cmdName)

	err := c.Execute(context.TODO(), func(ctx context.Context) error {
		return nil
	}, nil)
	assert.Nil(t, err)

	// this execute will open the circuit
	err = c.Execute(context.TODO(), func(ctx context.Context) error {
		return errRunFunc
	}, nil)
	assert.Equal(t, errRunFunc, err)
	m := new(dto.Metric)
	ot, _ := openedTotal.GetMetricWithLabelValues(cmdName)
	o, _ := opened.GetMetricWithLabelValues(cmdName)
	_ = ot.Write(m)
	_ = o.Write(m)
	assert.Equal(t, 1, int(*m.Counter.Value))
	assert.Equal(t, 1, int(*m.Gauge.Value))

	// wait the circuit close
	time.Sleep(1 * time.Second)

	_ = c.Execute(context.TODO(), func(ctx context.Context) error {
		return nil
	}, nil)
	m = new(dto.Metric)
	_ = o.Write(m)
	assert.Equal(t, 0, int(*m.Gauge.Value))

	err = c.Execute(context.TODO(), func(ctx context.Context) error {
		return errRunFunc
	}, func(ctx context.Context, err error) error {
		assert.Equal(t, errRunFunc, err)
		return errFallbackFunc
	})
	assert.Equal(t, errFallbackFunc, err)

	if os.Getenv("DEBUG_SLEEP") != "" {
		time.Sleep(1 * time.Hour)
	}
}


func TestCommandFactoryWithNameSpace(t *testing.T) {
	f := rolling.StatFactory{}
	prom := GetFactoryWithNameSpace(prometheus.DefaultRegisterer, "")
	h := circuit.Manager{
		DefaultCircuitProperties: []circuit.CommandPropertiesConstructor{
			f.CreateConfig,
			prom.CommandProperties,
			func(circuitName string) circuit.Config {
				return circuit.Config{
					Execution: circuit.ExecutionConfig{
						MaxConcurrentRequests: 10,
					},
					General: circuit.GeneralConfig{
						OpenToClosedFactory: hystrix.CloserFactory(hystrix.ConfigureCloser{
							SleepWindow: time.Millisecond * 800,
						}),
						ClosedToOpenFactory: hystrix.OpenerFactory(hystrix.ConfigureOpener{
							RequestVolumeThreshold:   2,
							ErrorThresholdPercentage: 50,
						}),
					},
				}
			},
		},
	}

	cmdName := "test-cmd"
	c := h.MustCreateCircuit(cmdName)

	err := c.Execute(context.TODO(), func(ctx context.Context) error {
		return nil
	}, nil)
	assert.Nil(t, err)

	// this execute will open the circuit
	err = c.Execute(context.TODO(), func(ctx context.Context) error {
		return errRunFunc
	}, nil)
	assert.Equal(t, errRunFunc, err)
	m := new(dto.Metric)
	ot, _ := openedTotal.GetMetricWithLabelValues(cmdName)
	o, _ := opened.GetMetricWithLabelValues(cmdName)
	_ = ot.Write(m)
	_ = o.Write(m)
	assert.Equal(t, 2, int(*m.Counter.Value))
	assert.Equal(t, 1, int(*m.Gauge.Value))

	// wait the circuit close
	time.Sleep(1 * time.Second)

	_ = c.Execute(context.TODO(), func(ctx context.Context) error {
		return nil
	}, nil)
	m = new(dto.Metric)
	_ = o.Write(m)
	assert.Equal(t, 0, int(*m.Gauge.Value))

	err = c.Execute(context.TODO(), func(ctx context.Context) error {
		return errRunFunc
	}, func(ctx context.Context, err error) error {
		assert.Equal(t, errRunFunc, err)
		return errFallbackFunc
	})
	assert.Equal(t, errFallbackFunc, err)

	if os.Getenv("DEBUG_SLEEP") != "" {
		time.Sleep(1 * time.Hour)
	}
}
