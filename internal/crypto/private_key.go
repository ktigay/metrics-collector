package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
)

// PrivateKey структура для хранения приватного ключа.
type PrivateKey struct {
	Key    *rsa.PrivateKey
	reader io.Reader
}

func (p *PrivateKey) init() error {
	var (
		err             error
		privateKeyPEM   []byte
		privateKeyBlock *pem.Block
		privateKey      *rsa.PrivateKey
	)

	if privateKeyPEM, err = io.ReadAll(p.reader); err != nil {
		return err
	}

	privateKeyBlock, _ = pem.Decode(privateKeyPEM)

	privateKey, err = x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return err
	}

	p.Key = privateKey
	return nil
}

// NewPrivateKey конструктор.
func NewPrivateKey(reader io.Reader) (*PrivateKey, error) {
	k := PrivateKey{
		reader: reader,
	}
	if err := k.init(); err != nil {
		return nil, err
	}

	return &k, nil
}
