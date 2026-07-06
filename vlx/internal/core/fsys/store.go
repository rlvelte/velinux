package fsys

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Store is a generic file-backed entity store.
type Store[T any] struct {
	Dir    string
	exts   []string
	decode DecodeFunc[T]
}

// DecodeFunc decodes a file into an entity.
type DecodeFunc[T any] func(name, path string, data []byte) (T, error)

// NewStore creates a Store backed by the given directory.
func NewStore[T any](baseDir string, decode DecodeFunc[T], exts ...string) *Store[T] {
	if len(exts) == 0 {
		exts = []string{""}
	}

	return &Store[T]{Dir: baseDir, exts: exts, decode: decode}
}

// List returns all decoded entities in the store.
func (s *Store[T]) List() ([]T, error) {
	entries, err := os.ReadDir(s.Dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var entities []T
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name, ok := s.matchExt(entry.Name())
		if !ok {
			continue
		}

		path := filepath.Join(s.Dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		entity, err := s.decode(name, path, data)
		if err != nil {
			continue
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

// Get retrieves a single entity by name.
func (s *Store[T]) Get(name string) (T, error) {
	for _, ext := range s.exts {
		path := filepath.Join(s.Dir, name+ext)
		data, err := os.ReadFile(path)

		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}

			var zero T
			return zero, err
		}

		return s.decode(name, path, data)
	}

	var zero T
	return zero, fmt.Errorf("%s not found", name)
}

// matchExt checks if the extension is in the specified field.
func (s *Store[T]) matchExt(filename string) (string, bool) {
	for _, ext := range s.exts {
		if ext == "" {
			return filename, true
		}

		if strings.HasSuffix(filename, ext) {
			name := strings.TrimSuffix(filename, ext)
			if name != "" {
				return name, true
			}
		}
	}

	return "", false
}
