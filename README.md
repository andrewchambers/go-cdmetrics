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
  -metrics-address string
        Address to send metrics to. (default "localhost:25826")
  -metrics-authfile string
        Path to metrics auth file. (default "/etc/collectd.authfile")
  -metrics-mode string
        Metrics mode, one of "disable", "plain-text", "sign", "encrypt". (default "plain-text")
  -metrics-username string
        Username for metrics auth modes. (default "metrics")

$ ./example
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
