package fsys

import (
	"os"
	"path/filepath"
)

// ConfigPath returns a path under XDG_CONFIG_HOME.
func ConfigPath(rel ...string) string {
	parts := append([]string{env("XDG_CONFIG_HOME", ".config")}, rel...)
	return filepath.Join(parts...)
}

// DataPath returns a path under XDG_DATA_HOME.
func DataPath(rel ...string) string {
	parts := append([]string{env("XDG_DATA_HOME", ".local/share")}, rel...)
	return filepath.Join(parts...)
}

// CachePath returns a path under XDG_CACHE_HOME.
func CachePath(rel ...string) string {
	parts := append([]string{env("XDG_CACHE_HOME", ".cache")}, rel...)
	return filepath.Join(parts...)
}

// env look for existing variable and construct fallback if not found.
func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, fallback)
}
