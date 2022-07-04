package cdmetrics

import (
	"sync"
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
	m := NewCounterMetric(name)
	c := &Counter{
		Metric: m,
		v:      0,
	}
	AddCollectorFunc(func(sink cdclient.MetricSink) error {
		return sink.AddValues(c.Metric, time.Now(), float64(c.Load()))
	})
	return c
}

type Gauge struct {
	Metric *cdclient.Metric
	m      sync.Mutex
	v      float64
}

func (g *Gauge) Load() float64 {
	g.m.Lock()
	defer g.m.Unlock()
	return g.v
}

func (g *Gauge) Set(v float64) {
	g.m.Lock()
	defer g.m.Unlock()
	g.v = v
}

func (g *Gauge) Add(v float64) {
	g.m.Lock()
	defer g.m.Unlock()
	g.v += v
}

func (g *Gauge) Inc() {
	g.Add(1.0)
}

func (g *Gauge) Dec() {
	g.Add(-1.0)
}

func NewGauge(name string) *Gauge {
	m := NewGaugeMetric(name)
	g := &Gauge{
		Metric: m,
		v:      0.0,
	}
	AddCollectorFunc(func(sink cdclient.MetricSink) error {
		return sink.AddValues(g.Metric, time.Now(), g.Load())
	})
	return g
}
