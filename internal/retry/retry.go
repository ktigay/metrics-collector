package retry

import (
	"time"
)

var defaultDelays = []time.Duration{
	time.Second,
	3 * time.Second,
	5 * time.Second,
}

// DelayHandler хендлер.
type DelayHandler func(policy RetPolicy) error

// RetPolicy политика ретраев.
type RetPolicy interface {
	Next() bool
	Retries() int
	LastError() error
	SetLastError(err error)
	SetSkip(bool)
	IsSkipped() bool
}

type defaultRetPolicy struct {
	delays    []time.Duration
	tries     int
	lastError error
	skip      bool
}

// NewDefaultRetPolicy конструктор.
func NewDefaultRetPolicy(delays []time.Duration) RetPolicy {
	return &defaultRetPolicy{delays: delays}
}

// Next true если можно делать следующую попытку вызова.
func (d *defaultRetPolicy) Next() bool {
	if d.lastError == nil || d.skip || d.tries >= len(d.delays) {
		return false
	}

	time.Sleep(d.delays[d.tries])
	d.tries++
	return true
}

// Retries кол-во повторных вызовов.
func (d *defaultRetPolicy) Retries() int {
	return d.tries
}

// SetLastError сохраняет последнюю ошибку.
func (d *defaultRetPolicy) SetLastError(err error) {
	d.lastError = err
}

// LastError возвращает последнюю ошибку.
func (d *defaultRetPolicy) LastError() error {
	return d.lastError
}

// SetSkip прервать ретраи.
func (d *defaultRetPolicy) SetSkip(skip bool) {
	d.skip = skip
}

// IsSkipped ретраи прерваны.
func (d *defaultRetPolicy) IsSkipped() bool {
	return d.skip
}

// RetWithDelays ретраи с политикой.
func RetWithDelays(handler DelayHandler, policy RetPolicy) error {
	var retryFn func() error

	retryFn = func() error {
		policy.SetLastError(
			handler(policy),
		)
		if !policy.Next() {
			return policy.LastError()
		}

		return retryFn()
	}

	return retryFn()
}

// Ret ретраи с дефолтной политикой.
func Ret(handler DelayHandler) error {
	return RetWithDelays(handler, NewDefaultRetPolicy(defaultDelays))
}
