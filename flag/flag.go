// cdflag can be imported to automatically add cdmetrics options for
// the standard library flag library.
package cdflag

import (
	"flag"

	"github.com/andrewchambers/go-cdmetrics"
)

func init() {
	flag.StringVar(&cdmetrics.MetricMode, "metrics-mode", cdmetrics.MetricMode, "Metrics mode, one of \"disabled\",\"unencrypted\", \"signed\", \"encrypted\".")
	flag.StringVar(&cdmetrics.MetricAddress, "metrics-address", cdmetrics.MetricAddress, "Address of collectd network address to send metrics to.")
	flag.StringVar(&cdmetrics.MetricUsername, "metrics-username", cdmetrics.MetricUsername, "Username for use sign and encrypt modes.")
	flag.StringVar(&cdmetrics.MetricAuthFile, "metrics-authfile", cdmetrics.MetricAuthFile, "Path to collectd auth file.")
}
