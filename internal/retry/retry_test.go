package retry

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetWithDelays(t *testing.T) {
	type args struct {
		handler DelayHandler
		delays  []time.Duration
	}
	tests := []struct {
		name      string
		args      args
		wantTries int
		wantErr   bool
	}{
		{
			name: "Fist_call_success",
			args: args{
				handler: func(_ RetPolicy) error {
					return nil
				},
				delays: []time.Duration{
					1 * time.Second,
				},
			},
			wantTries: 1,
			wantErr:   false,
		},
		{
			name: "Retries_with_error",
			args: args{
				handler: func(_ RetPolicy) error {
					return fmt.Errorf("some error")
				},
				delays: []time.Duration{
					10 * time.Millisecond,
					10 * time.Millisecond,
					10 * time.Millisecond,
					10 * time.Millisecond,
				},
			},
			wantTries: 5,
			wantErr:   true,
		},
		{
			name: "Retries_with_skip",
			args: args{
				handler: func(policy RetPolicy) error {
					if policy.Retries() == 1 {
						policy.SetSkip(true)
					}
					return fmt.Errorf("some error")
				},
				delays: []time.Duration{
					10 * time.Millisecond,
					10 * time.Millisecond,
					10 * time.Millisecond,
					10 * time.Millisecond,
				},
			},
			wantTries: 2,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tries int
			handler := func(policy RetPolicy) error {
				tries++
				return tt.args.handler(policy)
			}
			if err := RetWithDelays(handler, NewDefaultRetPolicy(tt.args.delays)); (err != nil) != tt.wantErr {
				t.Errorf("RetWithDelays() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.wantTries, tries)
		})
	}
}
