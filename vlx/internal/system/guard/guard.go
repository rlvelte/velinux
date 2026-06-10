package guard

import (
	"errors"
	"net"
	"os"
	"strings"
	"time"
)

var (
	errUnsupportedPlatform = errors.New("unsupported platform")
	errNoConnection        = errors.New("no internet connection")
)

const osReleasePath = "/etc/os-release"

// OS enforces that the application runs on supported distributions.
func OS() error {
	info, err := os.ReadFile(osReleasePath)
	if err != nil {
		return errUnsupportedPlatform
	}

	content := strings.ToLower(string(info))
	if !strings.Contains(content, "opensuse") {
		return errUnsupportedPlatform
	}

	return nil
}

// Connection checks for basic internet connectivity.
func Connection() error {
	conn, err := net.DialTimeout("tcp", "1.1.1.1:53", 3*time.Second)
	if err != nil {
		return errNoConnection
	}

	_ = conn.Close()
	return nil
}
