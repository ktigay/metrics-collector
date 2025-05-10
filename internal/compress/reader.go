package compress

import (
	"compress/gzip"
	"compress/zlib"
	"io"

	"github.com/andybalholm/brotli"
)

// Reader структура для чтения сжатых данных.
type Reader struct {
	reader io.ReadCloser
	cmp    io.ReadCloser
}

// Read распаковка сжатых данных.
func (c *Reader) Read(p []byte) (n int, err error) {
	return c.cmp.Read(p)
}

// Close закрытие ридера.
func (c *Reader) Close() error {
	if err := c.reader.Close(); err != nil {
		return err
	}
	return c.cmp.Close()
}

// ReaderFactory фабрика для создания ридера сжатых данных.
func ReaderFactory(t Type, r io.ReadCloser) (io.ReadCloser, error) {
	switch t {
	case Gzip:
		return newGZipCompressReader(r)
	case Deflate:
		return newDeflateCompressReader(r)
	case Br:
		return newBrotliCompressReader(r)
	default:
		return r, nil
	}
}

func newGZipCompressReader(r io.ReadCloser) (io.ReadCloser, error) {
	cmp, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &Reader{
		reader: r,
		cmp:    cmp,
	}, nil
}

func newDeflateCompressReader(r io.ReadCloser) (io.ReadCloser, error) {
	cmp, err := zlib.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &Reader{
		reader: r,
		cmp:    cmp,
	}, nil
}

func newBrotliCompressReader(r io.ReadCloser) (io.ReadCloser, error) {
	cmp := brotli.NewReader(r)

	return &Reader{
		reader: r,
		cmp: brotliDecorator{
			cmp,
		},
	}, nil
}

type brotliDecorator struct {
	*brotli.Reader
}

func (br brotliDecorator) Close() error {
	return nil
}
