package snapshot

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/ktigay/metrics-collector/internal"
	"github.com/ktigay/metrics-collector/internal/server/storage"
)

type FileSnapshot struct {
	filePath string
}

func NewFileSnapshot(filePath string) *FileSnapshot {
	return &FileSnapshot{filePath: filePath}
}

func (f *FileSnapshot) Read() ([]storage.Entity, error) {
	return FileReadAll[storage.Entity](f.filePath)
}

func (f *FileSnapshot) Write(entities []storage.Entity) error {
	return FileWriteAll[storage.Entity](f.filePath, entities)
}

// FileWriteAll запись структур в виде json-строк в файл.
func FileWriteAll[T any](path string, e []T) error {
	if err := ensureDir(filepath.Dir(path)); err != nil {
		return err
	}

	writer, err := NewAtomicFileWriter(path)
	if err != nil {
		return err
	}

	for _, el := range e {
		if err = writer.Write(el); err != nil {
			return err
		}
	}

	if err = writer.Flush(); err != nil {
		return err
	}
	return writer.Close()
}

// FileReadAll чтение json-строк в структуры из файла.
func FileReadAll[T any](path string) ([]T, error) {
	if err := ensureDir(filepath.Dir(path)); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer internal.Quite(file.Close)
	var all = make([]T, 0)

	dec := json.NewDecoder(file)

	for {
		var e T
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

func tempDir(dest string) string {
	tmp := os.Getenv("TMPDIR")
	if tmp == "" {
		tmp = filepath.Dir(dest)
	}
	return tmp
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
