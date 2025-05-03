package io

import (
	"compress/gzip"
	"compress/zlib"
	"io"

	"github.com/andybalholm/brotli"
	"github.com/ktigay/metrics-collector/internal/compress"
)

// CompressReader структура для чтения сжатых данных.
type CompressReader struct {
	reader io.ReadCloser
	cmp    io.ReadCloser
}

// Read распаковка сжатых данных.
func (c *CompressReader) Read(p []byte) (n int, err error) {
	return c.cmp.Read(p)
}

// Close закрытие ридера.
func (c *CompressReader) Close() error {
	if err := c.reader.Close(); err != nil {
		return err
	}
	return c.cmp.Close()
}

// CompressReaderFactory фабрика для создания ридера сжатых данных.
func CompressReaderFactory(t compress.Type, r io.ReadCloser) (io.ReadCloser, error) {
	switch t {
	case compress.Gzip:
		return newGZipCompressReader(r)
	case compress.Deflate:
		return newDeflateCompressReader(r)
	case compress.Br:
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

	return &CompressReader{
		reader: r,
		cmp:    cmp,
	}, nil
}

func newDeflateCompressReader(r io.ReadCloser) (io.ReadCloser, error) {
	cmp, err := zlib.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &CompressReader{
		reader: r,
		cmp:    cmp,
	}, nil
}

func newBrotliCompressReader(r io.ReadCloser) (io.ReadCloser, error) {
	cmp := brotli.NewReader(r)

	return &CompressReader{
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
