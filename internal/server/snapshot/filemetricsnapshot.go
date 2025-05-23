package snapshot

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/ktigay/metrics-collector/internal/log"
	"github.com/ktigay/metrics-collector/internal/server/storage"
	"go.uber.org/zap"
)

// FileMetricSnapshot структура для сохранения снапшота в файле.
type FileMetricSnapshot struct {
	filePath string
}

// NewFileMetricSnapshot конструктор.
func NewFileMetricSnapshot(filePath string) *FileMetricSnapshot {
	return &FileMetricSnapshot{filePath: filePath}
}

// Read чтение снапшота из файла.
func (f *FileMetricSnapshot) Read() ([]storage.MetricEntity, error) {
	if err := ensureDir(filepath.Dir(f.filePath)); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(f.filePath, os.O_RDONLY|os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.AppLogger.Error("snapshot.Read error", zap.Error(err))
		}
	}()
	all := make([]storage.MetricEntity, 0)

	dec := json.NewDecoder(file)

	for {
		var e storage.MetricEntity
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
func (f *FileMetricSnapshot) Write(entities []storage.MetricEntity) error {
	if err := ensureDir(filepath.Dir(f.filePath)); err != nil {
		return err
	}

	writer, err := NewAtomicFileWriter(f.filePath)
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
