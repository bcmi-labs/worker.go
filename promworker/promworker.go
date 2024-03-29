/*
* Copyright 2019 ARDUINO SA (http://www.arduino.cc/)
* This file is part of worker.go.
* Copyright (c) 2019
* Authors: Matteo "triex" Suppo
*
* This software is released under:
* first, the GNU General Public License, which covers the main part of
* worker.go
* The terms of this license can be found at:
* https://www.gnu.org/licenses/gpl-3.0.en.html
*
* next, the GNU Lesser General Public License that covers
* worker.go. The terms of this license can be found at:
* https://www.gnu.org/licenses/lgpl.html
*
* You can be released from the requirements of the above licenses by purchasing
* a commercial license. Buying such a license is mandatory if you want to modify or
* otherwise use the software for commercial activities involving the Arduino
* software without disclosing the source code of your own applications. To purchase
* a commercial license, send an email to license@arduino.cc.
*
 */

package promworker

import (
	"time"

	"github.com/bcmi-labs/worker.go"
	"github.com/prometheus/client_golang/prometheus"
)

// New creates a worker.Pool backed by prometheus stats.
// It will create and register:
// - a {{prefix}}_jobs_count gauge keeping track of the number of jobs currently running
// - a {{prefix}}_jobs_sum counter keeping track of the total number of jobs ran
// - a {{prefix}}_trace_nanoseconds summary keeping track of the duration time of traces
// You have to start the prometheus server yourself, though
func New(prefix string) *worker.Pool {
	count := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: prefix + "_jobs_count",
		Help: "The number of jobs currently running in the worker pool",
	})
	prometheus.MustRegister(count)

	sum := prometheus.NewCounter(prometheus.CounterOpts{
		Name: prefix + "_jobs_sum",
		Help: "The number of jobs currently running in the worker pool",
	})
	prometheus.MustRegister(sum)

	jobs := jobs{
		count: count,
		sum:   sum,
	}

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
		Jobs:   &jobs,
		Tracer: &tracer,
	}
}

type jobs struct {
	count prometheus.Gauge
	sum   prometheus.Counter
}

func (j *jobs) Inc() {
	j.count.Inc()
	j.sum.Inc()
}

func (j *jobs) Dec() {
	j.count.Dec()
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
