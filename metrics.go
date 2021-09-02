// An opinionated abstraction for collectd metrics.
package cdmetrics

import (
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/andrewchambers/go-cdclient"
)

var (
	LogFn                = log.Printf
	MetricInterval       = 10 * time.Second
	MetricHost           string
	MetricPlugin         string = "gocdmetrics"
	MetricPluginInstance string

	UDPAddress  string = "localhost:25826"
	UDPMode     string = "disabled"
	UDPUsername string = "metrics"
	UDPAuthFile string = "/etc/collectd.authfile"
)

var (
	mlock      sync.Mutex
	collectors []func(cdclient.MetricSink) error
)

func init() {
	host, _ := os.Hostname()
	if host == "" {
		host = "_unknown_"
	}
	MetricHost = host
	MetricPluginInstance = path.Base(os.Args[0])
}

func Start() {
	if UDPMode == "disabled" {
		return
	}
	go udpMetricsForever()
}

func CollectInto(sink cdclient.MetricSink) error {
	mlock.Lock()
	defer mlock.Unlock()
	for _, f := range collectors {
		err := f(sink)
		if err != nil {
			return err
		}
	}
	return nil
}

func udpMetricsForever() {

	opts := cdclient.UDPClientOptions{
		BufferSize: cdclient.DefaultBufferSize,
	}

	switch UDPMode {
	case "unencrypted":
		opts.Mode = cdclient.UDPPlainText
	case "signed":
		opts.Mode = cdclient.UDPSign
		opts.Username = UDPUsername
	case "encrypted":
		opts.Mode = cdclient.UDPEncrypt
		opts.Username = UDPUsername
	default:
		LogFn("invalid metrics mode %q, metrics disabled", UDPMode)
		return
	}

	if UDPMode != "unencrypted" {
		authFile, err := cdclient.NewAuthFile(UDPAuthFile)
		if err != nil {
			LogFn("unable to load %q: %s", UDPAuthFile, err)
			return
		}
		password, ok := authFile.Password(UDPUsername)
		if !ok {
			LogFn("no password for -metrics-user %q", UDPUsername)
			return
		}
		opts.Password = password
	}

	for {
		client, err := cdclient.DialUDP(UDPAddress, opts)
		if err != nil {
			LogFn("unable to create metrics client: %s", err)
			time.Sleep(MetricInterval)
			continue
		}

		for {

			err = CollectInto(client)
			if err != nil {
				LogFn("metrics collection failed: %s", err)
				break
			}

			err = client.Flush()
			if err != nil {
				LogFn("metrics client flush failed: %s", err)
				break
			}
			time.Sleep(MetricInterval)
		}
	}

}

func NewDefaultMetric() *cdclient.Metric {
	return &cdclient.Metric{
		Host:           MetricHost,
		Plugin:         MetricPlugin,
		PluginInstance: MetricPluginInstance,
		Interval:       MetricInterval,
	}
}

func AddCollectorFunc(f func(cdclient.MetricSink) error) {
	mlock.Lock()
	defer mlock.Unlock()
	collectors = append(collectors, f)
}
