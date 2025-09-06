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
	return &FileSnapper[T]{filename: filename}
}

func (f *FileSnapper[T]) Snap(_ context.Context, data T) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(f.filename, d, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileSnapper[T]) Load(_ context.Context) (T, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var data T
	cont, err := os.ReadFile(f.filename)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(cont, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
