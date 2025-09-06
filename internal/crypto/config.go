package crypto

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

// Config конфигурация.
type Config struct {
	PrivateKeyPath string `env:"PRIVATE_KEY_PATH"`
	PublicKeyPath  string `env:"PUBLIC_KEY_PATH"`
}

const (
	defaultPrivateKeyPath = "./certs/private.key"
	defaultPublicKeyPath  = "./certs/certificate.pem"
)

// InitializeConfig инициализирует конфиг клиента.
func InitializeConfig(args []string) (*Config, error) {
	config := Config{
		PrivateKeyPath: defaultPrivateKeyPath,
		PublicKeyPath:  defaultPublicKeyPath,
	}

	flags := flag.NewFlagSet("crypto flags", flag.ContinueOnError)

	flags.StringVar(&config.PrivateKeyPath, "k", defaultPrivateKeyPath, "private key path")
	flags.StringVar(&config.PublicKeyPath, "s", defaultPublicKeyPath, "certificate path")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
