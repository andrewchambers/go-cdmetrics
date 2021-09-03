package main

import (
	"flag"
	"time"

	"github.com/andrewchambers/go-cdmetrics"
	"github.com/andrewchambers/go-cdmetrics/rtmetrics"

	_ "github.com/andrewchambers/go-cdmetrics/flag"
)


func main() {

	flag.Parse()

	cdmetrics.MetricPluginInstance = "example"
	counter := cdmetrics.NewCounter("my-counter")
	rtmetrics.RegisterGoRuntimeMetrics()

	cdmetrics.Start()

	for {
		counter.Inc()
		time.Sleep(1)
	}

}
