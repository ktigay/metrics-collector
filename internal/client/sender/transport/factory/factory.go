package factory

import (
	"bufio"
	"log"
	"net/http"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ktigay/metrics-collector/internal/client/config"
	"github.com/ktigay/metrics-collector/internal/client/sender"
	"github.com/ktigay/metrics-collector/internal/client/sender/transport"
	"github.com/ktigay/metrics-collector/internal/contracts"
	"github.com/ktigay/metrics-collector/internal/crypto"
)

// CreateTransport создает транспорт для передачи данных на сервер.
func CreateTransport(cfg *config.Config, logger *zap.SugaredLogger) sender.Transport {
	switch cfg.Transport {
	case config.TransportGRPC:
		return getGRPCTransport(cfg, logger)
	default:
		return getHTTPTransport(cfg, logger)
	}
}

func getHTTPTransport(
	cfg *config.Config,
	logger *zap.SugaredLogger,
) sender.Transport {
	var (
		cryptoKey *crypto.PublicKey
		err       error
	)
	if cfg.CryptoKey != "" {
		if cryptoKey, err = initPublicKey(cfg.CryptoKey); err != nil {
			log.Fatalf("can't initialize crypto key: %v", err)
		}
	}

	factory := transport.NewRequestFactory(http.MethodPost, cfg.ServerProtocol+"://"+cfg.ServerHost, cfg.HashKey, cfg.IPAddr, cryptoKey)
	return transport.NewHTTPClient(factory, logger)
}

func getGRPCTransport(cfg *config.Config, logger *zap.SugaredLogger) sender.Transport {
	conn, err := grpc.NewClient(cfg.ServerGRPCHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("can't create gRPC transport: %v", err)
	}

	return transport.NewGRPCClient(contracts.NewMetricsServiceClient(conn), logger)
}

func initPublicKey(cryptoKey string) (*crypto.PublicKey, error) {
	file, err := os.Open(cryptoKey)
	if err != nil {
		return nil, err
	}
	return crypto.NewPublicKey(bufio.NewReader(file))
}
