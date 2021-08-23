// cdflag can be imported to automatically add cdmetrics options for
// the standard library flag library.
package cdflag

import (
	"flag"

	"github.com/andrewchambers/go-cdmetrics"
)

func init() {
	flag.StringVar(&cdmetrics.MetricAddress, "metrics-address", cdmetrics.MetricAddress, "Address to send metrics to.")
	flag.StringVar(&cdmetrics.MetricUsername, "metrics-username", cdmetrics.MetricUsername, "Username for metrics auth modes.")
	flag.StringVar(&cdmetrics.MetricMode, "metrics-mode", cdmetrics.MetricMode, "Metrics mode, one of \"disable\", \"plain-text\", \"sign\", \"encrypt\".")
	flag.StringVar(&cdmetrics.MetricAuthFile, "metrics-authfile", cdmetrics.MetricAuthFile, "Path to metrics auth file.")
}
