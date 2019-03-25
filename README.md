worker
======

Worker is a worker pool thought to be easily observable. You can plug it with whatever you want, but we provide out of the box a worker pool
integrated with [Prometheus](https://prometheus.io/)

# Usage

```go
package main

import (
	"log"
	"net/http"
	"time"

	// Include these!
	"github.com/bcmi-labs/worker.go"
	"github.com/bcmi-labs/worker.go/promworker"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Start the prometheus handler!
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Listening prometheus metrics on %s/metrics\n", ":8888")
	go func() {
		log.Fatal(http.ListenAndServe(":8888", nil))
	}()

	// Instantiate a new worker with the name test!
	pool := promworker.New("test")

	// Run a job that sleeps for 10 seconds!
	pool.Run(func() {
		time.Sleep(time.Second * 10)
	})

	// Run another job that registers the duration of actions inside it!
	pool.RunTraced("trace", func(trace worker.Trace) {
		time.Sleep(time.Second * 10)
		// When you call Mark it's like you stopped a stopwatch and recorded
		// the time from the last Mark or the start of the trace!
		trace.Mark("first action")
		time.Sleep(time.Second * 5)
		trace.Mark("second action")
		time.Sleep(time.Second * 1)
		trace.Mark("third action")
	})

	// Block forever!
	select {}
}

```

This program will expose a prometheus endpoint with the following metrics:

```
# HELP test_jobs_count The number of jobs currently running in the worker pool
# TYPE test_jobs_count gauge
test_jobs_count 2
```

```
# HELP test_trace_nanoseconds The duration in nanoseconds of arbitrary actions performed by jobs in the worker pool
# TYPE test_trace_nanoseconds summary
test_trace_nanoseconds{action="first action",scope="trace",quantile="0.5"} 1.0000253174e+10
test_trace_nanoseconds{action="first action",scope="trace",quantile="0.9"} 1.0000253174e+10
test_trace_nanoseconds{action="first action",scope="trace",quantile="0.99"} 1.0000253174e+10
test_trace_nanoseconds_sum{action="first action",scope="trace"} 1.0000253174e+10
test_trace_nanoseconds_count{action="first action",scope="trace"} 1
test_trace_nanoseconds{action="second action",scope="trace",quantile="0.5"} 5.000191237e+09
test_trace_nanoseconds{action="second action",scope="trace",quantile="0.9"} 5.000191237e+09
test_trace_nanoseconds{action="second action",scope="trace",quantile="0.99"} 5.000191237e+09
test_trace_nanoseconds_sum{action="second action",scope="trace"} 5.000191237e+09
test_trace_nanoseconds_count{action="second action",scope="trace"} 1
test_trace_nanoseconds{action="third action",scope="trace",quantile="0.5"} 1.000101584e+09
test_trace_nanoseconds{action="third action",scope="trace",quantile="0.9"} 1.000101584e+09
test_trace_nanoseconds{action="third action",scope="trace",quantile="0.99"} 1.000101584e+09
test_trace_nanoseconds_sum{action="third action",scope="trace"} 1.000101584e+09
test_trace_nanoseconds_count{action="third action",scope="trace"} 1
```