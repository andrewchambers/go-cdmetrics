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
	MetricPublishInterval = 10 * time.Second
	MetricHost            string
	MetricPlugin          string = "go"
	MetricPluginInstance  string

	MetricAddress  string = "localhost:25826"
	MetricUsername string = "metrics"
	MetricMode     string = "plain-text"
	MetricAuthFile string = "/etc/collectd.authfile"
)

var (
	mlock      sync.Mutex
	collectors []func(*cdclient.UDPClient) error
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
	go metricsForever()
}

func metricsForever() {

	opts := cdclient.UDPClientOptions{
		BufferSize: cdclient.DefaultBufferSize,
	}

	switch MetricMode {
	case "disable":
		return
	case "plain-text":
		opts.Mode = cdclient.UDPPlainText
	case "sign":
		opts.Mode = cdclient.UDPSign
		opts.Username = MetricUsername
	case "encrypt":
		opts.Mode = cdclient.UDPEncrypt
		opts.Username = MetricUsername
	default:
		log.Printf("invalid -metrics-mode %q", MetricMode)
		return
	}

	if MetricMode != "plain-text" {
		authFile, err := cdclient.NewAuthFile(MetricAuthFile)
		if err != nil {
			log.Printf("unable to load %q: %s", MetricAuthFile, err)
			return
		}
		password, ok := authFile.Password(MetricUsername)
		if !ok {
			log.Printf("no password for -metrics-user %q", MetricUsername)
			return
		}
		opts.Password = password
	}

	for {
		client, err := cdclient.DialUDP(MetricAddress, opts)
		if err != nil {
			log.Printf("unable to create metrics client: %s", err)
			time.Sleep(MetricPublishInterval)
			continue
		}

		for {
			mlock.Lock()
			for _, f := range collectors {
				err := f(client)
				if err != nil {
					log.Printf("metric collector failed: %s", err)
				}
			}
			mlock.Unlock()

			err = client.Flush()
			if err != nil {
				log.Printf("metrics client flush failed: %s", err)
				break
			}
			time.Sleep(MetricPublishInterval)
		}
	}

}

func NewDefaultMetric() *cdclient.Metric {
	return &cdclient.Metric{
		Host:           MetricHost,
		Plugin:         MetricPlugin,
		PluginInstance: MetricPluginInstance,
		Interval:       MetricPublishInterval,
	}
}

func AddCollectorFunc(f func(*cdclient.UDPClient) error) {
	mlock.Lock()
	defer mlock.Unlock()
	collectors = append(collectors, f)
}
