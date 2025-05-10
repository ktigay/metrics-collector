package snapshot

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
)

type Encoder interface {
	Encode(interface{}) error
}

type AtomicFileWriter struct {
	tmpFile  *os.File
	filePath string
	writer   *bufio.Writer
	encoder  Encoder
}

func NewAtomicFileWriter(filePath string) (*AtomicFileWriter, error) {
	tmp, err := os.CreateTemp(tempDir(filePath), "atomic-*")
	if err != nil {
		return nil, err
	}
	writer := bufio.NewWriter(tmp)
	return &AtomicFileWriter{
		filePath: filePath,
		tmpFile:  tmp,
		writer:   writer,
		encoder:  defaultEncoder(writer),
	}, nil
}

// Write запись данных.
func (a *AtomicFileWriter) Write(e any) error {
	err := a.encoder.Encode(e)
	if err != nil {
		a.onError()
	}

	return err
}

// Flush записывает данные из буфера в источник.
func (a *AtomicFileWriter) Flush() error {
	err := a.writer.Flush()
	if err != nil {
		a.onError()
	}

	return err
}

func (a *AtomicFileWriter) Close() error {
	err := a.tmpFile.Chmod(0644)
	defer func() {
		if err != nil {
			a.onError()
		}
	}()

	if err != nil {
		return err
	}

	err = a.tmpFile.Sync()
	if err != nil {
		return err
	}

	err = a.tmpFile.Close()
	if err != nil {
		return err
	}

	err = os.Rename(a.tmpFile.Name(), a.filePath)
	if err != nil {
		a.onError()
	}

	return err
}

func (a *AtomicFileWriter) onError() {
	if a.tmpFile != nil {
		_ = a.tmpFile.Close()
		_ = os.Remove(a.tmpFile.Name())
	}
}

func defaultEncoder(w io.Writer) Encoder {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc
}
