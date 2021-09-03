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
	m.DSTypes = []cdclient.DSType{cdclient.COUNTER}
	c := &Counter{Metric: m}
	AddCollectorFunc(func(sink cdclient.MetricSink) error {
		return sink.AddValues(c.Metric, time.Now(), float64(c.Load()))
	})
	return c
}

type Gauge struct {
	Metric *cdclient.Metric
	v      uint64
}

func (g *Gauge) Load() float64 {
	return math.Float64frombits(atomic.LoadUint64(&g.v))
}

func (g *Gauge) Set(v float64) {
	atomic.StoreUint64(&g.v, math.Float64bits(v))
}

func (g *Gauge) Add(v float64) {
	for {
		old := atomic.LoadUint64(&g.v)
		new := math.Float64bits(math.Float64frombits(old) + 1)
		if atomic.CompareAndSwapUint64(&g.v, old, new) {
			break
		}
	}
}

func (g *Gauge) Inc() {
	g.Add(1.0)
}

func (g *Gauge) Dec() {
	g.Add(-1.0)
}

func NewGauge(name string) *Gauge {
	m := NewDefaultMetric()
	m.Type = "gauge"
	m.TypeInstance = name
	m.DSTypes = []cdclient.DSType{cdclient.GAUGE}
	g := &Gauge{Metric: m}
	g.Set(0.0)
	AddCollectorFunc(func(sink cdclient.MetricSink) error {
		return sink.AddValues(g.Metric, time.Now(), g.Load())
	})
	return g
}
