package compress

import "testing"

func TestTypeFromString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want Type
	}{
		{
			name: "Positive_test_#1",
			args: args{
				str: "gzip",
			},
			want: Gzip,
		},
		{
			name: "Positive_test_#2",
			args: args{
				str: "deflate,  br",
			},
			want: Deflate,
		},
		{
			name: "Positive_test_#3",
			args: args{
				str: "br,deflate",
			},
			want: Br,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TypeFromString(tt.args.str); got != tt.want {
				t.Errorf("TypeFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}
