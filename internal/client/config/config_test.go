package config

import (
	"os"
	"reflect"
	"testing"
)

func Test_parseFlags(t *testing.T) {
	type args struct {
		filePath     string
		fileContents string
		envs         map[string]string
		flags        []string
	}
	tests := []struct {
		name    string
		want    *Config
		args    args
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
				CryptoKey:      defaultCryptoKey,
				RateLimit:      defaultRateLimit,
				ServerGRPCHost: defaultGRPCHost,
				Transport:      defaultTransport,
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
				CryptoKey:      defaultCryptoKey,
				RateLimit:      defaultRateLimit,
				ServerGRPCHost: defaultGRPCHost,
				Transport:      defaultTransport,
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
				CryptoKey:      defaultCryptoKey,
				RateLimit:      defaultRateLimit,
				ServerGRPCHost: defaultGRPCHost,
				Transport:      defaultTransport,
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
				CryptoKey:      defaultCryptoKey,
				RateLimit:      defaultRateLimit,
				ServerGRPCHost: defaultGRPCHost,
				Transport:      defaultTransport,
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
		{
			name: "Positive_test_Envs_Flags_JSON",
			args: args{
				envs: map[string]string{
					"REPORT_INTERVAL": "111",
				},
				flags:    []string{"-r=120", "-c=/tmp/agent.json"},
				filePath: "/tmp/agent.json",
				fileContents: `{
				  "address": "localhost:8088",
				  "report_interval": "1s",
				  "poll_interval": "1s",
				  "crypto_key": "./certs/crypto.pem"
				}`,
			},
			want: &Config{
				ServerProtocol: defaultServerProtocol,
				ServerHost:     "localhost:8088",
				ConfigFile:     "/tmp/agent.json",
				ReportInterval: 111,
				PollInterval:   1,
				LogLevel:       defaultLogLevel,
				CryptoKey:      "./certs/crypto.pem",
				RateLimit:      defaultRateLimit,
				ServerGRPCHost: defaultGRPCHost,
				Transport:      defaultTransport,
			},
			wantErr: false,
		},
		{
			name: "Positive_test_Envs_Flags_JSON_From_Args_With_Space",
			args: args{
				envs: map[string]string{
					"REPORT_INTERVAL": "111",
				},
				flags:    []string{"-r=120", "-c", "/tmp/agent.json", "-p=11"},
				filePath: "/tmp/agent.json",
				fileContents: `{
				  "address": "localhost:8088",
				  "report_interval": "1s",
				  "poll_interval": "1s",
				  "crypto_key": "./certs/crypto.pem"
				}`,
			},
			want: &Config{
				ServerProtocol: defaultServerProtocol,
				ServerHost:     "localhost:8088",
				ConfigFile:     "/tmp/agent.json",
				ReportInterval: 111,
				PollInterval:   11,
				LogLevel:       defaultLogLevel,
				CryptoKey:      "./certs/crypto.pem",
				RateLimit:      defaultRateLimit,
				ServerGRPCHost: defaultGRPCHost,
				Transport:      defaultTransport,
			},
			wantErr: false,
		},
		{
			name: "Positive_test_Envs_Flags_JSON_From_Env",
			args: args{
				envs: map[string]string{
					"REPORT_INTERVAL": "111",
					"CONFIG":          "/tmp/agent.json",
				},
				flags:    []string{"-r=120"},
				filePath: "/tmp/agent.json",
				fileContents: `{
				  "address": "localhost:8088",
				  "report_interval": "1s",
				  "poll_interval": "1s",
				  "crypto_key": "./certs/crypto.pem"
				}`,
			},
			want: &Config{
				ServerProtocol: defaultServerProtocol,
				ServerHost:     "localhost:8088",
				ConfigFile:     "/tmp/agent.json",
				ReportInterval: 111,
				PollInterval:   1,
				LogLevel:       defaultLogLevel,
				CryptoKey:      "./certs/crypto.pem",
				RateLimit:      defaultRateLimit,
				ServerGRPCHost: defaultGRPCHost,
				Transport:      defaultTransport,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			if tt.args.filePath != "" && tt.args.fileContents != "" {
				f, err := os.OpenFile(tt.args.filePath, os.O_WRONLY|os.O_CREATE, 0o644)
				if err != nil {
					panic(err)
				}
				defer func() {
					_ = f.Close()
					_ = os.Remove(tt.args.filePath)
				}()

				if _, err = f.WriteString(tt.args.fileContents); err != nil {
					panic(err)
				}
			}
			if tt.args.envs != nil {
				for k, v := range tt.args.envs {
					if err := os.Setenv(k, v); err != nil {
						t.Fatal(err)
					}
				}
			}
			got, err := NewConfig(tt.args.flags)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
