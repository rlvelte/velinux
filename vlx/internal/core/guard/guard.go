package guard

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

const osReleasePath = "/etc/os-release"

// OS enforces that the application runs on supported distributions.
func OS() error {
	info, err := os.ReadFile(osReleasePath)
	if err != nil {
		return err
	}

	content := strings.ToLower(string(info))
	if !strings.Contains(content, "opensuse") {
		return fmt.Errorf("%s is not an OpenSUSE release", osReleasePath)
	}

	return nil
}

// Connection checks for basic internet connectivity.
func Connection() error {
	conn, err := net.DialTimeout("tcp", "1.1.1.1:53", 3*time.Second)
	if err != nil {
		return fmt.Errorf("you seem to be offline: %s", err)
	}

	_ = conn.Close()
	return nil
}

// Binaries verifies that all required executables are available on PATH.
func Binaries(required ...string) error {
	for _, bin := range required {
		if _, err := exec.LookPath(bin); err != nil {
			return errors.New("required binary not found: " + bin)
		}
	}
	return nil
}
