package worker

type Pool struct {
	Jobs interface {
		Inc()
		Dec()
	}
}

func (p *Pool) Run(fn func()) {
	p.Jobs.Inc()

	go func() {
		fn()

		p.Jobs.Dec()
	}()
}
