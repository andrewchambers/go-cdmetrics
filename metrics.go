// An opinionated abstraction for collectd metrics.
package cdmetrics

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/andrewchambers/go-cdclient"
)

var (
	LogFn                 = log.Printf
	MetricsInterval       = 10 * time.Second
	MetricsHost           string
	MetricsPlugin         string = "gocdmetrics"
	MetricsPluginInstance string
	MetricsMode           string = "disabled"

	UDPAddress  string = "localhost:25826"
	UDPUsername string = "metrics"
	UDPAuthFile string = "/etc/collectd.authfile"
)

var (
	mlock      sync.Mutex
	metrics    []*cdclient.Metric
	collectors []func(cdclient.MetricSink) error
)

func init() {
	host, _ := os.Hostname()
	if host == "" {
		host = "_unknown_"
	}
	MetricsHost = host
}

func Start() {

	if MetricsMode == "disabled" {
		return
	}

	// We must configure metrics that were created after
	// Global values were configured.
	mlock.Lock()
	defer mlock.Unlock()
	for _, m := range metrics {
		m.Interval = MetricsInterval
		m.Host = MetricsHost
		m.Plugin = MetricsPlugin
		m.PluginInstance = MetricsPluginInstance
		err := m.Validate()
		if err != nil {
			panic("error validating metric: " + err.Error())
		}
	}

	if strings.HasSuffix(MetricsMode, "-udp") {
		go udpMetricsForever()
	} else {
		LogFn("invalid metrics mode %q, metrics disabled", MetricsMode)
	}
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

	switch MetricsMode {
	case "unencrypted-udp":
		opts.Mode = cdclient.UDPPlainText
	case "signed-udp":
		opts.Mode = cdclient.UDPSign
		opts.Username = UDPUsername
	case "encrypted-udp":
		opts.Mode = cdclient.UDPEncrypt
		opts.Username = UDPUsername
	default:
		LogFn("invalid udp metrics mode %q, metrics disabled", MetricsMode)
		return
	}

	if MetricsMode != "unencrypted-udp" {
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
			time.Sleep(MetricsInterval)
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
			time.Sleep(MetricsInterval)
		}
	}

}

func NewMetric() *cdclient.Metric {
	m := &cdclient.Metric{
		Host:           MetricsHost,
		Plugin:         MetricsPlugin,
		PluginInstance: MetricsPluginInstance,
		Interval:       MetricsInterval,
	}
	metrics = append(metrics, m)
	return m
}

var guageDSTypes []cdclient.DSType = []cdclient.DSType{cdclient.GAUGE}

func NewGaugeMetric(instance string) *cdclient.Metric {
	m := &cdclient.Metric{
		Host:           MetricsHost,
		Plugin:         MetricsPlugin,
		PluginInstance: MetricsPluginInstance,
		Interval:       MetricsInterval,
		Type:           "gauge",
		TypeInstance:   instance,
		DSTypes:        guageDSTypes,
	}
	metrics = append(metrics, m)
	return m
}

var counterDSTypes []cdclient.DSType = []cdclient.DSType{cdclient.COUNTER}

func NewCounterMetric(instance string) *cdclient.Metric {
	m := &cdclient.Metric{
		Host:           MetricsHost,
		Plugin:         MetricsPlugin,
		PluginInstance: MetricsPluginInstance,
		Interval:       MetricsInterval,
		Type:           "counter",
		TypeInstance:   instance,
		DSTypes:        counterDSTypes,
	}
	metrics = append(metrics, m)
	return m
}

func AddCollectorFunc(f func(cdclient.MetricSink) error) {
	mlock.Lock()
	defer mlock.Unlock()
	collectors = append(collectors, f)
}
