package fsys

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// DecodeFunc decodes a file into an entity.
type DecodeFunc[T any] func(name, path string, data []byte) (T, error)

// GetJSON reads a single JSON file from dir/name.json and decodes it.
func GetJSON[T any](dir, name string, decode DecodeFunc[T]) (T, error) {
	path := filepath.Join(dir, name+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		var zero T
		return zero, err
	}

	return decode(name, path, data)
}

// ListJSON reads and decodes all .json files in a directory.
func ListJSON[T any](dir string, decode DecodeFunc[T]) ([]T, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}

		return nil, err
	}

	var entities []T
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".json")
		if name == "" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		entity, err := decode(name, path, data)
		if err != nil {
			continue
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

// AtomicWrite writes data to a path atomically using a temporary file + rename.
func AtomicWrite(path string, data []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmp, path)
}

// SetJSON marshals v and writes it atomically to dir/name.json.
// The directory is created if it does not exist.
func SetJSON[T any](dir, name string, v T) error {
	path := filepath.Join(dir, name+".json")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return AtomicWrite(path, data)
}
