package client

import (
	"os"
	"reflect"
	"testing"
)

func Test_parseFlags(t *testing.T) {
	type args struct {
		envs  map[string]string
		flags []string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "Positive_test_Default_Values",
			want: &Config{
				ServerProtocol: defaultServerProtocol,
				ServerHost:     defaultServerHost,
				ReportInterval: defaultReportInterval,
				PollInterval:   defaultPollInterval,
				LogLevel:       defaultLogLevel,
				RateLimit:      defaultRateLimit,
			},
			wantErr: false,
		},
		{
			name: "Positive_test_Envs",
			args: args{
				envs: map[string]string{
					"ADDRESS":         "localhost:8090",
					"REPORT_INTERVAL": "100",
					"POLL_INTERVAL":   "8",
				},
			},
			want: &Config{
				ServerProtocol: defaultServerProtocol,
				ServerHost:     "localhost:8090",
				ReportInterval: 100,
				PollInterval:   8,
				LogLevel:       defaultLogLevel,
				RateLimit:      defaultRateLimit,
			},
			wantErr: false,
		},
		{
			name: "Positive_test_Flags",
			args: args{
				envs: map[string]string{
					"ADDRESS":         "",
					"REPORT_INTERVAL": "",
					"POLL_INTERVAL":   "",
				},
				flags: []string{"-a=localhost:80100", "-r=120", "-p=15"},
			},
			want: &Config{
				ServerProtocol: defaultServerProtocol,
				ServerHost:     "localhost:80100",
				ReportInterval: 120,
				PollInterval:   15,
				LogLevel:       defaultLogLevel,
				RateLimit:      defaultRateLimit,
			},
			wantErr: false,
		},
		{
			name: "Positive_test_Envs_Flags",
			args: args{
				envs: map[string]string{
					"ADDRESS":         "localhost:8099",
					"REPORT_INTERVAL": "111",
					"POLL_INTERVAL":   "7",
				},
				flags: []string{"-a=localhost:80100", "-r=120", "-p=15"},
			},
			want: &Config{
				ServerProtocol: defaultServerProtocol,
				ServerHost:     "localhost:8099",
				ReportInterval: 111,
				PollInterval:   7,
				LogLevel:       defaultLogLevel,
				RateLimit:      defaultRateLimit,
			},
			wantErr: false,
		},
		{
			name: "Negative_test_Address_Invalid",
			args: args{
				envs: map[string]string{
					"ADDRESS":         "",
					"REPORT_INTERVAL": "",
					"POLL_INTERVAL":   "",
				},
				flags: []string{"-a=", "-r=120", "-p=15"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Negative_test_Address_Invalid_white_space",
			args: args{
				envs: map[string]string{
					"ADDRESS":         " ",
					"REPORT_INTERVAL": "122",
					"POLL_INTERVAL":   "2",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Negative_test_Report_Interval_Invalid_empty",
			args: args{
				envs: map[string]string{
					"ADDRESS":         "",
					"REPORT_INTERVAL": "",
					"POLL_INTERVAL":   "",
				},
				flags: []string{"-a=localhost:80100", "-r=0", "-p=15"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Negative_test_Report_Interval_Invalid_zero_value",
			args: args{
				envs: map[string]string{
					"ADDRESS":         "",
					"REPORT_INTERVAL": "0",
					"POLL_INTERVAL":   "",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Negative_test_Poll_Interval_Invalid_empty",
			args: args{
				envs: map[string]string{
					"ADDRESS":         "",
					"REPORT_INTERVAL": "",
					"POLL_INTERVAL":   "",
				},
				flags: []string{"-a=localhost:80100", "-r=122", "-p=0"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Negative_test_Poll_Interval_Invalid_zero_value",
			args: args{
				envs: map[string]string{
					"ADDRESS":         "",
					"REPORT_INTERVAL": "",
					"POLL_INTERVAL":   "0",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Negative_test_Invalid_Flag",
			args: args{
				envs: map[string]string{
					"ADDRESS":         "",
					"REPORT_INTERVAL": "",
					"POLL_INTERVAL":   "",
				},
				flags: []string{"-a=localhost:80100", "-r=122", "-p=11", "-s=12"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.envs != nil {
				for k, v := range tt.args.envs {
					if err := os.Setenv(k, v); err != nil {
						t.Fatal(err)
					}
				}
			}
			got, err := InitializeConfig(tt.args.flags)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitializeConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitializeConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
