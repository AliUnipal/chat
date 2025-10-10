package snapper

import (
	"context"
	"encoding/json"
	"os"
	"sync"
)

type Snapper[T any] interface {
	Snap(ctx context.Context, data T) error
	Load(ctx context.Context) (T, error)
}

type FileSnapper[T any] struct {
	filename string
	mu       sync.RWMutex
}

func NewFileSnapper[T any](filename string) *FileSnapper[T] {
	// TODO: filename and path validation
	return &FileSnapper[T]{filename: filename}
}

func (f *FileSnapper[T]) Snap(_ context.Context, data T) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	file, err := os.OpenFile(f.filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(data)
}

func (f *FileSnapper[T]) Load(_ context.Context) (T, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var data T
	file, err := os.Open(f.filename)
	if err != nil {
		return data, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return data, err
	}

	return data, nil
}
