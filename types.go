package cdmetrics

import (
	"math"
	"sync/atomic"
	"time"

	"github.com/andrewchambers/go-cdclient"
)

type Counter struct {
	Metric *cdclient.Metric
	v      int64
}

func (ctr *Counter) Inc() {
	atomic.AddInt64(&ctr.v, 1)
}

func (ctr *Counter) Add(i int64) {
	if i > 0 {
		atomic.AddInt64(&ctr.v, i)
	}
}

func (ctr *Counter) Load() int64 {
	return atomic.LoadInt64(&ctr.v)
}

func NewCounter(name string) *Counter {
	m := NewDefaultMetric()
	m.Type = "counter"
	m.TypeInstance = name
	m.ValueTypes = []cdclient.ValueType{cdclient.COUNTER}
	m.Validate()
	c := &Counter{Metric: m}
	AddCollectorFunc(func(client *cdclient.UDPClient) error {
		return client.AddValues(c.Metric, time.Now(), float64(c.Load()))
	})
	return c
}

type Gauge struct {
	Metric *cdclient.Metric
	v      uint64
}

func (g *Gauge) Set(i float64) {
	atomic.StoreUint64(&g.v, math.Float64bits(i))
}

func (g *Gauge) Load() float64 {
	return math.Float64frombits(atomic.LoadUint64(&g.v))
}

func NewGauge(name string) *Gauge {
	m := NewDefaultMetric()
	m.Type = "gauge"
	m.TypeInstance = name
	m.ValueTypes = []cdclient.ValueType{cdclient.GAUGE}
	m.Validate()
	g := &Gauge{Metric: m}
	g.Set(0.0)
	AddCollectorFunc(func(client *cdclient.UDPClient) error {
		return client.AddValues(g.Metric, time.Now(), g.Load())
	})
	return g
}
