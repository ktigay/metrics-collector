package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"io"
)

// Writer имплементация Writer для шифрования.
type Writer struct {
	Writer io.Writer
	Key    *rsa.PublicKey
}

// Write шифрует данные в RSA PKCS1v15.
func (e *Writer) Write(p []byte) (n int, err error) {
	var b []byte
	b, err = rsa.EncryptPKCS1v15(rand.Reader, e.Key, p)
	if err != nil {
		return 0, err
	}
	return e.Writer.Write(b)
}
