package main

import (
	"flag"
	"time"

	"github.com/andrewchambers/go-cdmetrics"
	_ "github.com/andrewchambers/go-cdmetrics/flag"
)

var counter = cdmetrics.NewCounter("my-counter")

func main() {
	flag.Parse()

	cdmetrics.ExportGoRuntimeMetrics()
	cdmetrics.Start()

	for {
		counter.Inc()
		time.Sleep(1)
	}

}
