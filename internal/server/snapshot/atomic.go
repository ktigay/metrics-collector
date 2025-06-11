package snapshot

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

// Encoder интерфейс энкодера.
type Encoder interface {
	Encode(interface{}) error
}

// AtomicFileWriter структура для атомарной записи в файл.
type AtomicFileWriter struct {
	tmpFile  *os.File
	filePath string
	writer   *bufio.Writer
	encoder  Encoder
	logger   *zap.SugaredLogger
}

// NewAtomicFileWriter конструктор.
func NewAtomicFileWriter(filePath string, logger *zap.SugaredLogger) (*AtomicFileWriter, error) {
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
		logger:   logger,
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

// Close закрываем запись.
func (a *AtomicFileWriter) Close() (err error) {
	defer func() {
		if err != nil {
			a.onError()
		}
	}()

	if err = a.tmpFile.Chmod(0o644); err != nil {
		return err
	}

	if err = a.tmpFile.Sync(); err != nil {
		return err
	}

	if err = a.tmpFile.Close(); err != nil {
		return err
	}

	err = os.Rename(a.tmpFile.Name(), a.filePath)

	return err
}

func (a *AtomicFileWriter) onError() {
	if a.tmpFile != nil {
		if err := a.tmpFile.Close(); err != nil {
			a.logger.Errorf("failed to close tmp file: %v", err)
		}
		if err := os.Remove(a.tmpFile.Name()); err != nil {
			a.logger.Errorf("failed to remove tmp file: %v", err)
		}
	}
}

func defaultEncoder(w io.Writer) Encoder {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc
}

func tempDir(dest string) string {
	tmp := os.Getenv("TMPDIR")
	if tmp == "" {
		tmp = filepath.Dir(dest)
	}
	return tmp
}
