// cdflag can be imported to automatically add cdmetrics options for
// the standard library flag library.
package cdflag

import (
	"flag"

	"github.com/andrewchambers/go-cdmetrics"
)

func init() {
	flag.StringVar(&cdmetrics.UDPMode, "udp-metrics-mode", cdmetrics.UDPMode, "Metrics mode, one of \"disabled\",\"unencrypted\", \"signed\", \"encrypted\".")
	flag.StringVar(&cdmetrics.UDPAddress, "udp-metrics-address", cdmetrics.UDPAddress, "Address of collectd network address to send metrics to.")
	flag.StringVar(&cdmetrics.UDPUsername, "udp-metrics-username", cdmetrics.UDPUsername, "Username for use sign and encrypt modes.")
	flag.StringVar(&cdmetrics.UDPAuthFile, "udp-metrics-authfile", cdmetrics.UDPAuthFile, "Path to collectd auth file.")
}
