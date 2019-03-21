package worker

type Pool struct {
	Jobs interface {
		Inc()
		Dec()
	}
	Recover func(interface{})
}

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
