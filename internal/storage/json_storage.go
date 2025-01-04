package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type JsonStorage[T any] struct {
	data     T
	filePath string
}

func (s *JsonStorage[T]) Read() error {
	if _, err := os.Stat(s.filePath); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	jsonData, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err := json.Unmarshal(jsonData, &s.data); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return nil
}

func (s *JsonStorage[T]) Write(data T) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	if err := os.WriteFile(s.filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

func (s *JsonStorage[T]) GetData() T {
	return s.data
}

func NewJsonStorage[T any](filePath string, initialData T) *JsonStorage[T] {
	return &JsonStorage[T]{filePath: filePath, data: initialData}
}
