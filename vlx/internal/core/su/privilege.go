package su

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rlvelte/velinux/vlx/internal/core/logs"
)

// RunPrivileged executes a command with privilege escalation unconditionally.
func RunPrivileged(args ...string) error {
	if err := run("sudo", append([]string{"-n"}, args...)...); err == nil {
		return nil
	}

	if path, err := exec.LookPath("vlxpass"); err == nil {
		sudoArgs := append([]string{"-A", "--preserve-env=WAYLAND_DISPLAY,HOME,XDG_RUNTIME_DIR"}, args...)
		cmd := exec.Command("sudo", sudoArgs...)
		cmd.Env = append(os.Environ(), "SUDO_ASKPASS="+path)

		cmd.Stdout = logs.Stdout()
		cmd.Stderr = logs.Stderr()

		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	if _, err := exec.LookPath("pkexec"); err == nil {
		return run("pkexec", args...)
	}

	return fmt.Errorf("privilege: no escalation method available")
}

func run(bin string, args ...string) error {
	cmd := exec.Command(bin, args...)
	cmd.Stdout = logs.Stdout()
	cmd.Stderr = logs.Stderr()
	return cmd.Run()
}
