package worker

// Pool is a worker pool thought to be easily observable
type Pool struct {
	// Jobs keep tracks of the number of jobs.
	// Whenever a job is started, Inc is called. Whenever a job finishes, Dec is called
	Jobs interface {
		Inc()
		Dec()
	}
	// Recover is an optional function that's launched to recover a panicing goroutine
	Recover func(interface{})
	// Tracer is used by RunTraced to spawn traces to keep track of the duration of actions during the running job
	Tracer Tracer
}

// Run spawns a new goroutine with the given function.
// Jobs.Inc() is called whenever it's enqueued, Jobs.Dec() is called when it finishes.
// Recover is called whenever the goroutine panics
func (p *Pool) Run(fn func()) {
	if p.Jobs != nil {
		p.Jobs.Inc()
	}

	go func() {
		defer func() {
			e := recover()
			if e != nil && p.Recover != nil {
				p.Recover(e)
			}
		}()
		fn()

		if p.Jobs != nil {
			p.Jobs.Dec()
		}
	}()
}

// RunTraced behaves like Run, but provides to the running goroutine a trace to keep track of the duration of actions
// Of course, if you don't provide a Tracer to the Pool, it will panic
func (p *Pool) RunTraced(scope string, fn func(trace Trace)) {
	if p.Jobs != nil {
		p.Jobs.Inc()
	}

	trace := p.Tracer.New(scope)

	go func() {
		defer func() {
			e := recover()
			if e != nil && p.Recover != nil {
				p.Recover(e)
			}
		}()
		fn(trace)

		if p.Jobs != nil {
			p.Jobs.Dec()
		}
	}()
}
