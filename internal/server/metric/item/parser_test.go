package item

import (
	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseFromPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantM   MetricDTO
		wantErr error
	}{
		{
			name: "Positive test",
			args: args{
				path: "/update/counter/pollCounter/12",
			},
			wantM: MetricDTO{
				Type:       metric.TypeCounter,
				Name:       "pollCounter",
				IntValue:   12,
				FloatValue: 0,
			},
			wantErr: nil,
		},
		{
			name: "Positive test #2",
			args: args{
				path: "/update/gauge/Alloc/22.01",
			},
			wantM: MetricDTO{
				Type:       metric.TypeGauge,
				Name:       "Alloc",
				IntValue:   0,
				FloatValue: 22.01,
			},
			wantErr: nil,
		},
		{
			name: "Error Invalid Name",
			args: args{
				path: "/update/counter/",
			},
			wantM:   MetricDTO{},
			wantErr: ErrorInvalidName,
		},
		{
			name: "Error Invalid Name #2",
			args: args{
				path: "/update/counter//122222.22",
			},
			wantM:   MetricDTO{},
			wantErr: ErrorInvalidName,
		},
		{
			name: "Error Invalid Length",
			args: args{
				path: "/update/gauge/Alloc/1222.000/111",
			},
			wantM:   MetricDTO{},
			wantErr: ErrorInvalidLength,
		},
		{
			name: "Error Invalid Length #2",
			args: args{
				path: "/update",
			},
			wantM:   MetricDTO{},
			wantErr: ErrorInvalidLength,
		},
		{
			name: "Error Invalid Type",
			args: args{
				path: "/update/wrongType/Alloc/1222.000",
			},
			wantM:   MetricDTO{},
			wantErr: ErrorInvalidType,
		},
		{
			name: "Error Invalid Value type",
			args: args{
				path: "/update/gauge/Alloc/floatValue",
			},
			wantM:   MetricDTO{},
			wantErr: ErrorInvalidVal,
		},
		{
			name: "Error Invalid Value type #2",
			args: args{
				path: "/update/counter/pollCount/intValue",
			},
			wantM:   MetricDTO{},
			wantErr: ErrorInvalidVal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotM, err := ParseFromPath(tt.args.path)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			}

			assert.Equal(t, tt.wantM, gotM)
		})
	}
}
