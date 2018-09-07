package metrics

import (
	"sync/atomic"
	"time"
)

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

type Counter struct {
	c uint64
}

// Value atomically returns the current counter value
func (c Counter) Value() uint64 {
	return atomic.LoadUint64(&c.c)
}

// Increment atomically increments the counter
func (c Counter) Increment() uint64 {
	return atomic.AddUint64(&c.c, 1)
}
