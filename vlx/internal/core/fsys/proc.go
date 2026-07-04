package fsys

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Proc reads a file from /proc.
func Proc(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join("/proc", name))
}

// ReadInt64 reads an int64 from a sysfs file.
func ReadInt64(path string) (int64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
}

// ReadUInt64 reads a uint64 from a sysfs file.
func ReadUInt64(path string) (uint64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
}
