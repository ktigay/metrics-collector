package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_IsUseSQLDB(t *testing.T) {
	type args struct {
		envs map[string]string
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "TestConfig_IsUseSQLDB_with_defaults",
			args: args{
				args: []string{},
				envs: map[string]string{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "TestConfig_IsUseSQLDB_with_envs",
			args: args{
				args: []string{},
				envs: map[string]string{
					"DATABASE_DSN": "postgres://postgres:postgres@localhost:2002/postgres?sslmode=disable",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "TestConfig_IsUseSQLDB_with_args",
			args: args{
				args: []string{
					"-d=postgres://postgres:postgres@localhost:1001/postgres?sslmode=disable",
				},
				envs: map[string]string{},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.args.envs {
				_ = os.Setenv(k, v)
			}
			got, err := NewConfig(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitializeConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got.IsUseSQLDB())
		})
	}
}

func TestNewConfig(t *testing.T) {
	type args struct {
		filePath     string
		fileContents string
		envs         map[string]string
		args         []string
	}
	tests := []struct {
		name    string
		want    *Config
		args    args
		wantErr bool
	}{
		{
			name: "TestInitializeConfig_with_flags",
			args: args{
				args: []string{
					"-a=localhost:1001",
					"-lvl=info",
					"-i=123",
					"-f=/tmp/restore-args.txt",
					"-r=0",
					"-d=postgres://postgres:postgres@localhost:1001/postgres?sslmode=disable",
				},
			},
			want: &Config{
				ServerHost:      "localhost:1001",
				LogLevel:        "info",
				StoreInterval:   123,
				FileStoragePath: "/tmp/restore-args.txt",
				Restore:         false,
				CryptoKey:       "",
				DatabaseDSN:     "postgres://postgres:postgres@localhost:1001/postgres?sslmode=disable",
				DatabaseDriver:  "pgx",
				ServerGRPCHost:  ":3000",
			},
		},
		{
			name: "TestInitializeConfig_with_args_and_envs",
			args: args{
				args: []string{
					"-a=localhost:1001",
					"-lvl=info",
					"-i=123",
					"-f=/tmp/restore-args.txt",
					"-r=0",
					"-d=postgres://postgres:postgres@localhost:1001/postgres?sslmode=disable",
				},
				envs: map[string]string{
					"ADDRESS":           "192.168.1.1:8080",
					"LOG_LEVEL":         "error",
					"STORE_INTERVAL":    "200",
					"FILE_STORAGE_PATH": "/tmp/restore-env.txt",
					"RESTORE":           "1",
					"DATABASE_DSN":      "postgres://postgres:postgres@localhost:2002/postgres?sslmode=disable",
					"DATABASE_DRIVER":   "mysql",
				},
			},
			want: &Config{
				ServerHost:      "192.168.1.1:8080",
				LogLevel:        "error",
				StoreInterval:   200,
				FileStoragePath: "/tmp/restore-env.txt",
				Restore:         true,
				CryptoKey:       "",
				DatabaseDSN:     "postgres://postgres:postgres@localhost:2002/postgres?sslmode=disable",
				DatabaseDriver:  "mysql",
				ServerGRPCHost:  ":3000",
			},
		},
		{
			name: "TestInitializeConfig_with_envs",
			args: args{
				args: []string{},
				envs: map[string]string{
					"ADDRESS":           "192.168.1.1:8080",
					"LOG_LEVEL":         "error",
					"STORE_INTERVAL":    "200",
					"FILE_STORAGE_PATH": "/tmp/restore-env.txt",
					"RESTORE":           "1",
					"DATABASE_DSN":      "postgres://postgres:postgres@localhost:2002/postgres?sslmode=disable",
					"DATABASE_DRIVER":   "mysql",
				},
			},
			want: &Config{
				ServerHost:      "192.168.1.1:8080",
				LogLevel:        "error",
				StoreInterval:   200,
				FileStoragePath: "/tmp/restore-env.txt",
				Restore:         true,
				CryptoKey:       "",
				DatabaseDSN:     "postgres://postgres:postgres@localhost:2002/postgres?sslmode=disable",
				DatabaseDriver:  "mysql",
				ServerGRPCHost:  ":3000",
			},
		},
		{
			name: "TestInitializeConfig_with_args_envs_and_json_file_from_args",
			args: args{
				args: []string{
					"-lvl=info",
					"-f=/tmp/restore-args.txt",
					"-r=0",
					"-d=postgres://postgres:postgres@localhost:1001/postgres?sslmode=disable",
					"-c=/tmp/config.json",
				},
				envs: map[string]string{
					"LOG_LEVEL":         "error",
					"FILE_STORAGE_PATH": "/tmp/restore-env.txt",
					"RESTORE":           "1",
					"DATABASE_DSN":      "postgres://postgres:postgres@localhost:2002/postgres?sslmode=disable",
					"DATABASE_DRIVER":   "mysql",
				},
				filePath: "/tmp/config.json",
				fileContents: `{
				  "address": ":80100",
				  "store_interval": "160s"
				}`,
			},
			want: &Config{
				ServerHost:      ":80100",
				LogLevel:        "error",
				ConfigFile:      "/tmp/config.json",
				StoreInterval:   160,
				FileStoragePath: "/tmp/restore-env.txt",
				Restore:         true,
				CryptoKey:       "",
				DatabaseDSN:     "postgres://postgres:postgres@localhost:2002/postgres?sslmode=disable",
				DatabaseDriver:  "mysql",
				ServerGRPCHost:  ":3000",
			},
		},
		{
			name: "TestInitializeConfig_with_args_envs_and_json_file_from_args_with_space_delimiter",
			args: args{
				args: []string{
					"-lvl=info",
					"-f=/tmp/restore-args.txt",
					"-r=0",
					"-c",
					"/tmp/config.json",
					"-d=postgres://postgres:postgres@localhost:1001/postgres?sslmode=disable",
				},
				envs: map[string]string{
					"LOG_LEVEL":         "error",
					"FILE_STORAGE_PATH": "/tmp/restore-env.txt",
					"RESTORE":           "1",
					"DATABASE_DSN":      "postgres://postgres:postgres@localhost:2002/postgres?sslmode=disable",
					"DATABASE_DRIVER":   "mysql",
				},
				filePath: "/tmp/config.json",
				fileContents: `{
				  "address": ":80100",
				  "store_interval": "160s"
				}`,
			},
			want: &Config{
				ServerHost:      ":80100",
				LogLevel:        "error",
				ConfigFile:      "/tmp/config.json",
				StoreInterval:   160,
				FileStoragePath: "/tmp/restore-env.txt",
				Restore:         true,
				CryptoKey:       "",
				DatabaseDSN:     "postgres://postgres:postgres@localhost:2002/postgres?sslmode=disable",
				DatabaseDriver:  "mysql",
				ServerGRPCHost:  ":3000",
			},
		},
		{
			name: "TestInitializeConfig_with_args_envs_and_json_file_from_env",
			args: args{
				args: []string{
					"-lvl=info",
					"-f=/tmp/restore-args.txt",
					"-r=0",
					"-d=postgres://postgres:postgres@localhost:1001/postgres?sslmode=disable",
				},
				envs: map[string]string{
					"LOG_LEVEL":         "error",
					"FILE_STORAGE_PATH": "/tmp/restore-env.txt",
					"RESTORE":           "1",
					"DATABASE_DSN":      "postgres://postgres:postgres@localhost:2002/postgres?sslmode=disable",
					"DATABASE_DRIVER":   "mysql",
					"CONFIG":            "/tmp/config.json",
				},
				filePath: "/tmp/config.json",
				fileContents: `{
				  "address": ":80100",
				  "restore": false,
				  "store_interval": "160s"
				}`,
			},
			want: &Config{
				ServerHost:      ":80100",
				LogLevel:        "error",
				ConfigFile:      "/tmp/config.json",
				StoreInterval:   160,
				FileStoragePath: "/tmp/restore-env.txt",
				Restore:         true,
				CryptoKey:       "",
				DatabaseDSN:     "postgres://postgres:postgres@localhost:2002/postgres?sslmode=disable",
				DatabaseDriver:  "mysql",
				ServerGRPCHost:  ":3000",
			},
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
			for k, v := range tt.args.envs {
				_ = os.Setenv(k, v)
			}
			got, err := NewConfig(tt.args.args)
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
