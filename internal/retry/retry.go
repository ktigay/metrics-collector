package retry

import (
	"time"
)

var defaultDelays = []time.Duration{
	time.Second,
	3 * time.Second,
	5 * time.Second,
}

// DelayHandler хендлер. Прерывает ретраи, если возвращает true
type DelayHandler func(p Policy) bool

// Policy политика ретраев.
type Policy struct {
	maxRetries int
	delays     []time.Duration
	retries    int
}

func (p *Policy) next() bool {
	idx := p.retries
	p.retries++

	if p.retries >= p.maxRetries {
		return false
	}

	var d time.Duration
	l := len(p.delays)
	if l > idx {
		d = p.delays[idx]
	} else {
		d = p.delays[l-1]
	}

	time.Sleep(d)
	return true
}

// RetIndex индекс ретраев.
func (p *Policy) RetIndex() int {
	return p.retries
}

// NewPolicy конструктор.
func NewPolicy(opts []Options) *Policy {
	p := &Policy{
		delays: defaultDelays,
	}
	for _, opt := range opts {
		opt(p)
	}
	if p.maxRetries == 0 && len(p.delays) > 0 {
		p.maxRetries = len(p.delays) + 1
	}
	return p
}

// Options Опции для ретраев.
type Options func(*Policy)

// WithRetries кол-во ретраев.
func WithRetries(max int) Options {
	return func(o *Policy) {
		o.maxRetries = max
	}
}

// WithDelays таймауты между ретраями. Если установлены, то maxRetries = len(delays) + 1
func WithDelays(delays []time.Duration) Options {
	return func(o *Policy) {
		o.delays = delays
	}
}

// Ret ретраи.
func Ret(handler DelayHandler, opts ...Options) {
	var retryFn func()
	p := NewPolicy(opts)

	retryFn = func() {
		if handler(*p) || !p.next() {
			return
		}
		retryFn()
	}

	retryFn()
}
