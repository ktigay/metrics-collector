package metric

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolveType(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantM   Type
		wantErr bool
	}{
		{
			name: "Positive_test",
			args: args{
				s: "gauge",
			},
			wantM:   TypeGauge,
			wantErr: false,
		},
		{
			name: "Negative test",
			args: args{
				s: "random_string",
			},
			wantM:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotM, err := ResolveType(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantM, gotM)
		})
	}
}
