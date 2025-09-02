package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
)

// PublicKey структура для хранения публичного ключа.
type PublicKey struct {
	Key    *rsa.PublicKey
	reader io.Reader
}

func (p *PublicKey) init() error {
	var (
		err          error
		publicKeyPEM []byte
		publicKey    any
	)

	if publicKeyPEM, err = io.ReadAll(p.reader); err != nil {
		return err
	}

	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	if publicKey, err = x509.ParsePKIXPublicKey(publicKeyBlock.Bytes); err != nil {
		return err
	}

	p.Key = publicKey.(*rsa.PublicKey)
	return nil
}

// NewPublicKey конструктор.
func NewPublicKey(reader io.Reader) (*PublicKey, error) {
	k := PublicKey{
		reader: reader,
	}
	if err := k.init(); err != nil {
		return nil, err
	}

	return &k, nil
}
