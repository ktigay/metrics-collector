package snapshot

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/ktigay/metrics-collector/internal"
	"io"
	"os"
)

// FileRead чтение из файла json-строки и распаковка в объект.
func FileRead[T any](path string) (*T, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer internal.Quite(file.Close)

	var e T
	if err := json.NewDecoder(file).Decode(&e); err != nil {
		if err != io.EOF {
			return nil, err
		}
		return nil, nil
	}
	return &e, nil
}

// FileWrite запись объекта в файл в виде json-строки.
func FileWrite[T any](path string, e *T) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer internal.Quite(file.Close)

	if err := json.NewEncoder(file).Encode(e); err != nil {
		return fmt.Errorf("failed write event: %v", err)
	}
	return nil
}

// FileWriteAll запись структур в виде json-строк в файл.
func FileWriteAll[T any](path string, e []*T) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer internal.Quite(file.Close)

	b := bufio.NewWriter(file)
	for _, et := range e {
		if err := json.NewEncoder(b).Encode(et); err != nil {
			return err
		}
	}
	internal.Quite(b.Flush)

	return nil
}

// FileReadAll чтение json-строк в структуры из файла.
func FileReadAll[T any](path string) ([]T, error) {
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
