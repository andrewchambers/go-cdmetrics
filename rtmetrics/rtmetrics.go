package rtmetrics

import (
	"runtime"
	"time"

	"github.com/andrewchambers/go-cdclient"
	"github.com/andrewchambers/go-cdmetrics"
)

func RegisterGoRuntimeMetrics() {
	vtGauge := []cdclient.DSType{
		cdclient.GAUGE,
	}

	goroutines := cdmetrics.NewDefaultMetric()
	goroutines.Type = "gauge"
	goroutines.TypeInstance = "rt-goroutines"
	goroutines.DSTypes = vtGauge

	heap_alloc := cdmetrics.NewDefaultMetric()
	heap_alloc.Type = "gauge"
	heap_alloc.TypeInstance = "rt-memstats-heap-alloc"
	heap_alloc.DSTypes = vtGauge

	gc_cpu_fraction := cdmetrics.NewDefaultMetric()
	gc_cpu_fraction.Type = "gauge"
	gc_cpu_fraction.TypeInstance = "rt-memstats-gc-cpu-fraction"
	gc_cpu_fraction.DSTypes = vtGauge

	cdmetrics.AddCollectorFunc(func(sink cdclient.MetricSink) error {
		t := time.Now()
		memStats := runtime.MemStats{}
		runtime.ReadMemStats(&memStats)
		_ = sink.AddValues(goroutines, t, float64(runtime.NumGoroutine()))
		_ = sink.AddValues(heap_alloc, t, float64(memStats.HeapAlloc))
		_ = sink.AddValues(gc_cpu_fraction, t, memStats.GCCPUFraction)
		return nil
	})
}
