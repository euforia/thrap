package metrics

import "time"

type Runtime struct {
	start time.Time
	end   time.Time
}

func (rt *Runtime) Start() *Runtime {
	rt.start = time.Now()
	return rt
}

func (rt *Runtime) End() {
	rt.end = time.Now()
}

func (rt *Runtime) Duration(round time.Duration) time.Duration {
	return rt.end.Sub(rt.start).Round(round)
}
