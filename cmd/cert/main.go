package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	c "github.com/ktigay/metrics-collector/internal/crypto"
)

func main() {
	var (
		err            error
		privateKey     *rsa.PrivateKey
		publicKeyBytes []byte
		cfg            *c.Config
	)

	cfg, err = c.InitializeConfig(os.Args[1:])
	if err != nil {
		log.Fatalf("Error initializing config: %s", err)
	}

	if _, err = os.Stat(cfg.PrivateKeyPath); err == nil {
		_, err = os.Stat(cfg.PublicKeyPath)
	}
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("can't stat %s: %v", cfg.PrivateKeyPath, err)
	}

	if privateKey, err = rsa.GenerateKey(rand.Reader, 4096); err != nil {
		log.Fatalf("Error generating private key: %s", err)
	}

	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	err = os.WriteFile(cfg.PrivateKeyPath, privateKeyPEM, 0o644)
	if err != nil {
		log.Fatalf("Error writing private key: %s", err)
	}

	publicKeyBytes, err = x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Fatalf("Error marshalling public key: %s", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	err = os.WriteFile(cfg.PublicKeyPath, publicKeyPEM, 0o644)
	if err != nil {
		log.Fatalf("Error writing public key: %s", err)
	}
}
