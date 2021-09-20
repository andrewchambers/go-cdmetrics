# go-cdmetrics

[![godocs.io](http://godocs.io/github.com/andrewchambers/go-cdmetrics?status.svg)](http://godocs.io/github.com/andrewchambers/go-cdmetrics)

A high level collectd metrics library for go. The metrics
are automatically exported to localhost via the collectd network
protocol.

# Example

./example/main.go:
```
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
```

In one terminal:

```
$ cd example
$ go build
$ ./example --help
Usage of ./example:
  -metrics-mode string
        Metrics mode, one of disabled, unencrypted-udp, signed-udp, encrypted-udp. (default "disabled")
  -metrics-udp-address string
        Address of collectd network address to send metrics to. (default "localhost:25826")
  -metrics-udp-authfile string
        Path to collectd auth file. (default "/etc/collectd.authfile")
  -metrics-udp-username string
        Username for use sign and encrypt modes. (default "metrics")

$ ./example -metrics-mode=encrypted-udp
```

In another terminal:
```
$ cd example
$ collectd -C ./collectd.conf -f
...
black.go-example.gauge-rt-goroutines 2 1629711890
black.go-example.gauge-rt-memstats-heap-alloc 500416 1629711890
black.go-example.counter-my-counter 80292077 1629711890
black.go-example.counter-my-counter 107281565 1629711900

```
