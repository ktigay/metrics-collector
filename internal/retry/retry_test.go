package retry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetWithDelays(t *testing.T) {
	type args struct {
		handler DelayHandler
		options []Options
	}
	tests := []struct {
		name      string
		args      args
		wantTries int
		wantSleep time.Duration
	}{
		{
			name: "Fist_call_success",
			args: args{
				handler: func(_ Policy) bool {
					return true
				},
				options: []Options{
					WithDelays([]time.Duration{
						1 * time.Second,
						1 * time.Second,
						1 * time.Second,
					}),
				},
			},
			wantTries: 1,
			wantSleep: 0,
		},
		{
			name: "Retries_with_max_tries",
			args: args{
				handler: func(_ Policy) bool {
					return false
				},
				options: []Options{
					WithDelays([]time.Duration{
						10 * time.Millisecond,
						20 * time.Millisecond,
						30 * time.Millisecond,
						40 * time.Millisecond,
					}),
				},
			},
			wantTries: 5,
			wantSleep: 100 * time.Millisecond,
		},
		{
			name: "Retries_with_skip",
			args: args{
				handler: func(policy Policy) bool {
					return policy.RetIndex() == 1
				},
				options: []Options{
					WithDelays([]time.Duration{
						10 * time.Millisecond,
						20 * time.Millisecond,
						30 * time.Millisecond,
					}),
				},
			},
			wantTries: 2,
			wantSleep: 10 * time.Millisecond,
		},
		{
			name: "Retries_with_retries_count",
			args: args{
				handler: func(policy Policy) bool {
					return false
				},
				options: []Options{
					WithRetries(5),
					WithDelays([]time.Duration{
						10 * time.Millisecond,
						40 * time.Millisecond,
					}),
				},
			},
			wantTries: 5,
			wantSleep: 130 * time.Millisecond,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var execCount int
			handler := func(policy Policy) bool {
				execCount++
				return tt.args.handler(policy)
			}

			start := time.Now()
			Ret(handler, tt.args.options...)
			elapsed := time.Since(start)

			assert.Equal(t, tt.wantTries, execCount)
			assert.Greater(t, elapsed, tt.wantSleep)
		})
	}
}
