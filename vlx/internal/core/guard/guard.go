package guard

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"time"
)

// Network checks for basic internet connectivity.
func Network() error {
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
