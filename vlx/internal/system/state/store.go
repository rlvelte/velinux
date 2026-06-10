package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rvelte/vlx/internal/system/xdg"
)

// Store persists JSON-encoded state under XDG_DATA_HOME/vlx/<namespace>/state.json.
type Store struct {
	path string
}

// New creates a store scoped to the given namespace.
func New(namespace string) *Store {
	return &Store{
		path: xdg.DataPath("vlx", namespace, "state.json"),
	}
}

// Path returns the underlying file path.
func (s *Store) Path() string {
	return s.path
}

// Read deserializes the state file into a value of type T.
func Read[T any](s *Store) (T, error) {
	var v T

	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return v, nil
		}
		return v, fmt.Errorf("read state: %w", err)
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return v, fmt.Errorf("parse state: %w", err)
	}

	return v, nil
}

// Write serializes v as indented JSON and persists it.
func (s *Store) Write(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal state: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("create state dir: %w", err)
	}

	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return fmt.Errorf("write state: %w", err)
	}

	return nil
}
