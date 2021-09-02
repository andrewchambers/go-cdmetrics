package main

import (
	"flag"
	"time"

	"github.com/andrewchambers/go-cdmetrics"
	"github.com/andrewchambers/go-cdmetrics/rtmetrics"

	_ "github.com/andrewchambers/go-cdmetrics/flag"
)

var counter = cdmetrics.NewCounter("my-counter")

func main() {
	cdmetrics.MetricPluginInstance = "example"

	flag.Parse()

	rtmetrics.RegisterGoRuntimeMetrics()
	cdmetrics.Start()

	for {
		counter.Inc()
		time.Sleep(1)
	}

}
