package io

import (
	"compress/gzip"
	"io"
)

type CompressReader struct {
	r  io.ReadCloser
	zr io.ReadCloser
}

func (c *CompressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *CompressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func CompressReaderFactory(t string, r io.ReadCloser) (io.ReadCloser, error) {
	switch t {
	case "gzip":
		return newGZipCompressReader(r)
	default:
		return r, nil
	}
}

func newGZipCompressReader(r io.ReadCloser) (io.ReadCloser, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &CompressReader{
		r:  r,
		zr: zr,
	}, nil
}
