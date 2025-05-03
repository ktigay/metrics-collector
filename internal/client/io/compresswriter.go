package io

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"github.com/andybalholm/brotli"
	iio "github.com/ktigay/metrics-collector/internal/compress"
	"io"
)

type CompressWriter struct {
	cmp io.WriteCloser
}

func (w *CompressWriter) Write(b []byte) (int, error) {
	return w.cmp.Write(b)
}

func (w *CompressWriter) Close() error {
	return w.cmp.Close()
}

func NewCompressWriter(t iio.Type, bb *bytes.Buffer) (*CompressWriter, error) {
	cmp, err := writer(t, bb)
	if err != nil {
		return nil, err
	}

	return &CompressWriter{
		cmp: cmp,
	}, nil
}

func writer(t iio.Type, bb *bytes.Buffer) (io.WriteCloser, error) {
	switch t {
	case iio.Gzip:
		return gzip.NewWriter(bb), nil
	case iio.Deflate:
		return zlib.NewWriter(bb), nil
	case iio.Br:
		return brotli.NewWriter(bb), nil
	}
	return nil, fmt.Errorf("unsupported compress type: %v", t)
}
