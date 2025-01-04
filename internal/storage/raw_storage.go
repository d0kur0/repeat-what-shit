package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type RawStorage struct {
	data     []byte
	filePath string
}

func (r *RawStorage) Read() error {
	if _, err := os.Stat(r.filePath); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err := json.Unmarshal(data, &r.data); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return nil
}

func (r *RawStorage) Write(data []byte) error {
	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

func (r *RawStorage) GetData() []byte {
	return r.data
}

func NewRawStorage(filePath string) *RawStorage {
	return &RawStorage{filePath: filePath}
}
