package main

import (
	"sync"
	"time"
)

type Clock interface{ Now() time.Time }

type Limiter struct {
	mu     sync.Mutex
	clock  Clock
	rate   float64
	burst  int
	tokens float64
	last   time.Time
}

func NewLimiter(clock Clock, ratePerSec float64, burst int) *Limiter {
	l := &Limiter{
		clock: clock,
		rate:  ratePerSec,
		burst: burst,
	}
	if clock == nil || burst <= 0 {
		return l
	}
	l.tokens = float64(burst)
	l.last = clock.Now()
	return l
}

func (l *Limiter) AllowN(n int) bool {
	if l.clock == nil || l.burst <= 0 || n <= 0 || n > l.burst {
		return false
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.rate > 0 {
		now := l.clock.Now()
		if now.After(l.last) {
			elapsed := now.Sub(l.last)
			l.tokens += float64(elapsed) * l.rate / float64(time.Second)
			if l.tokens > float64(l.burst) {
				l.tokens = float64(l.burst)
			}
			l.last = now
		}
	}
	need := float64(n)
	if l.tokens < need {
		return false
	}
	l.tokens -= need
	return true
}
