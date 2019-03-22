package promworker

import (
	"time"

	"github.com/bcmi-labs/worker"
	"github.com/prometheus/client_golang/prometheus"
)

// New creates a worker.Pool backed by prometheus stats.
// It will create and register:
// - a {{prefix}}_jobs_count gauge keeping track of the number of jobs
// - a {{prefix}}_trace_nanoseconds summary keeping track of the duration time of traces
// You have to start the prometheus server yourself, though
func New(prefix string) *worker.Pool {
	jobs := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: prefix + "_jobs_count",
		Help: "The number of jobs currently running in the worker pool",
	})
	prometheus.MustRegister(jobs)

	traces := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       prefix + "_trace_nanoseconds",
		Help:       "The duration in nanoseconds of arbitrary actions performed by jobs in the worker pool",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"action", "scope"})
	prometheus.MustRegister(traces)

	tracer := Tracer{
		traces: traces,
	}

	return &worker.Pool{
		Jobs:   jobs,
		Tracer: &tracer,
	}
}

// Tracer keeps track of the duration of jobs launched with pool.RunTraced()
// with a prometheus summary with the following labels:
// - scope: used to identify traces
// - action: used to identify actions inside of traces
type Tracer struct {
	traces *prometheus.SummaryVec
}

// New creates a new Trace with the given scope
func (t *Tracer) New(scope string) worker.Trace {
	return &Trace{
		scope:  scope,
		time:   time.Now(),
		traces: t.traces,
	}
}

// Trace keeps track of the duration of actions
type Trace struct {
	scope  string
	time   time.Time
	traces *prometheus.SummaryVec
}

// Mark observes the time elapsed from the previous Mark or the creation of the Trace
func (t *Trace) Mark(action string) {
	elapsed := time.Since(t.time).Nanoseconds()
	t.traces.With(prometheus.Labels{
		"scope":  t.scope,
		"action": action,
	}).Observe(float64(elapsed))
	t.time = time.Now()
}
