package rtmetrics

import (
	"runtime"
	"time"

	"github.com/andrewchambers/go-cdclient"
	"github.com/andrewchambers/go-cdmetrics"
)

func init() {
	goroutines := cdmetrics.NewGaugeMetric("rt-goroutines")
	heap_alloc := cdmetrics.NewGaugeMetric("rt-memstats-heap-alloc")
	gc_cpu_fraction := cdmetrics.NewGaugeMetric("rt-memstats-gc-cpu-fraction")
	cdmetrics.AddCollectorFunc(func(sink cdclient.MetricSink) error {
		t := time.Now()
		memStats := runtime.MemStats{}
		runtime.ReadMemStats(&memStats)
		err := sink.AddValues(goroutines, t, float64(runtime.NumGoroutine()))
		if err != nil {
			return err
		}
		err = sink.AddValues(heap_alloc, t, float64(memStats.HeapAlloc))
		if err != nil {
			return err
		}
		err = sink.AddValues(gc_cpu_fraction, t, memStats.GCCPUFraction)
		if err != nil {
			return err
		}
		return nil
	})
}
