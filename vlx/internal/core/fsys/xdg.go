package fsys

import (
	"os"
	"path/filepath"
	"strings"
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

// DataDirs returns the system data directories from XDG_DATA_DIRS.
func DataDirs() []string {
	s := os.Getenv("XDG_DATA_DIRS")
	if s == "" {
		s = "/usr/local/share:/usr/share"
	}
	return strings.Split(s, ":")
}

// CachePath returns a path under XDG_CACHE_HOME.
func CachePath(rel ...string) string {
	parts := append([]string{env("XDG_CACHE_HOME", ".cache")}, rel...)
	return filepath.Join(parts...)
}

// StatePath returns a path under XDG_STATE_HOME.
func StatePath(rel ...string) string {
	parts := append([]string{env("XDG_STATE_HOME", ".local/state")}, rel...)
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
