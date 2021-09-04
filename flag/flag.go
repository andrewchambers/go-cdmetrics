// cdflag can be imported to automatically add cdmetrics options for
// the standard library flag library.
package cdflag

import (
	"flag"

	"github.com/andrewchambers/go-cdmetrics"
)

func init() {
	flag.StringVar(&cdmetrics.MetricsMode, "metrics-mode", cdmetrics.MetricsMode, "Metrics mode, one of disabled, unencrypted-udp, signed-udp, encrypted-udp.")
	flag.StringVar(&cdmetrics.UDPAddress, "metrics-udp-address", cdmetrics.UDPAddress, "Address of collectd network address to send metrics to.")
	flag.StringVar(&cdmetrics.UDPUsername, "metrics-udp-username", cdmetrics.UDPUsername, "Username for use sign and encrypt modes.")
	flag.StringVar(&cdmetrics.UDPAuthFile, "metrics-udp-authfile", cdmetrics.UDPAuthFile, "Path to collectd auth file.")
}
