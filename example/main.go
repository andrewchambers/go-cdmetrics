package main

import (
	"flag"
	"time"

	"github.com/andrewchambers/go-cdmetrics"

	_ "github.com/andrewchambers/go-cdmetrics/flag"
	_ "github.com/andrewchambers/go-cdmetrics/rtmetrics"
)

var (
	counter = cdmetrics.NewCounter("my-counter")
)

func main() {

	flag.Parse()

	cdmetrics.MetricsPluginInstance = "example"
	cdmetrics.Start()

	for {
		counter.Inc()
		time.Sleep(1)
	}

}
