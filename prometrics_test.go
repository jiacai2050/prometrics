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
						MaxConcurrentRequests: 1000,
					},
					General: circuit.GeneralConfig{
						OpenToClosedFactory: hystrix.CloserFactory(hystrix.ConfigureCloser{
							SleepWindow: time.Second * 3,
						}),
						ClosedToOpenFactory: hystrix.OpenerFactory(hystrix.ConfigureOpener{
							RequestVolumeThreshold:   6,
							ErrorThresholdPercentage: 50,
						}),
					},
				}
			},
		},
	}

	c := h.MustCreateCircuit("hello-world")

	err := c.Execute(context.TODO(), func(ctx context.Context) error {
		return nil
	}, nil)
	assert.Nil(t, err)

	err = c.Execute(context.TODO(), func(ctx context.Context) error {
		return errRunFunc
	}, nil)

	assert.Equal(t, errRunFunc, err)

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
