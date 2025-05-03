package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
)

// FileRead чтение из файла json-строки и распаковка в объект.
func FileRead[T any](path string) (*T, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	var e T
	if err := json.NewDecoder(file).Decode(&e); err != nil {
		return nil, nil
	}
	return &e, nil
}

// FileWrite запись объекта в файл в виде json-строки.
func FileWrite[T any](path string, e *T) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	if err := json.NewEncoder(file).Encode(e); err != nil {
		return fmt.Errorf("failed write event: %v", err)
	}
	return nil
}
