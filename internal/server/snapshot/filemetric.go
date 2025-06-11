// Package snapshot Работа со снапшотами.
package snapshot

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/ktigay/metrics-collector/internal/server/repository"
	"go.uber.org/zap"
)

// FileMetricSnapshot структура для сохранения снапшота в файле.
type FileMetricSnapshot struct {
	filePath string
	logger   *zap.SugaredLogger
}

// NewFileMetricSnapshot конструктор.
func NewFileMetricSnapshot(filePath string, logger *zap.SugaredLogger) *FileMetricSnapshot {
	return &FileMetricSnapshot{
		filePath: filePath,
		logger:   logger,
	}
}

// Read чтение снапшота из файла.
func (f *FileMetricSnapshot) Read() ([]repository.MetricEntity, error) {
	if err := ensureDir(filepath.Dir(f.filePath)); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(f.filePath, os.O_RDONLY|os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = file.Close(); err != nil {
			f.logger.Error("snapshot.Read error", zap.Error(err))
		}
	}()
	all := make([]repository.MetricEntity, 0)

	dec := json.NewDecoder(file)

	for {
		var e repository.MetricEntity
		if err = dec.Decode(&e); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		all = append(all, e)
	}

	return all, nil
}

// Write запись данных в файл.
func (f *FileMetricSnapshot) Write(entities []repository.MetricEntity) error {
	if err := ensureDir(filepath.Dir(f.filePath)); err != nil {
		return err
	}

	writer, err := NewAtomicFileWriter(f.filePath, f.logger)
	if err != nil {
		return err
	}

	for _, el := range entities {
		if err = writer.Write(el); err != nil {
			return err
		}
	}

	if err = writer.Flush(); err != nil {
		return err
	}
	// перезапись происходит только при успешном закрытии writer.
	return writer.Close()
}

func ensureDir(dirName string) error {
	_, err := os.Stat(dirName)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(dirName, os.ModeDir)
}
