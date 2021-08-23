package cdmetrics

import (
	"runtime"
	"time"

	"github.com/andrewchambers/go-cdclient"
)

func ExportGoRuntimeMetrics() {
	vtGauge := []cdclient.ValueType{
		cdclient.GAUGE,
	}

	goroutines := NewDefaultMetric()
	goroutines.Type = "gauge"
	goroutines.TypeInstance = "rt-goroutines"
	goroutines.ValueTypes = vtGauge
	goroutines.Validate()

	heap_alloc := NewDefaultMetric()
	heap_alloc.Type = "gauge"
	heap_alloc.TypeInstance = "rt-memstats-heap-alloc"
	heap_alloc.ValueTypes = vtGauge
	heap_alloc.Validate()

	gc_cpu_fraction := NewDefaultMetric()
	gc_cpu_fraction.Type = "gauge"
	gc_cpu_fraction.TypeInstance = "rt-memstats-gc-cpu-fraction"
	gc_cpu_fraction.ValueTypes = vtGauge
	gc_cpu_fraction.Validate()

	AddCollectorFunc(func(client *cdclient.UDPClient) error {
		t := time.Now()
		memStats := runtime.MemStats{}
		runtime.ReadMemStats(&memStats)
		_ = client.AddValues(goroutines, t, float64(runtime.NumGoroutine()))
		_ = client.AddValues(heap_alloc, t, float64(memStats.HeapAlloc))
		_ = client.AddValues(gc_cpu_fraction, t, memStats.GCCPUFraction)
		return nil
	})
}
