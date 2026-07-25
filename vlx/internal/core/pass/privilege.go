package pass

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/rlvelte/velinux/vlx/internal/core/logs"
)

const ContextKey = "escalation"

// RunContext executes a command with privilege escalation when escalation is enabled.
func RunContext(ctx context.Context, args ...string) error {
	if ctx.Value(ContextKey).(bool) {
		return Run(args...)
	}

	if len(args) == 0 {
		return fmt.Errorf("pass: no command provided")
	}

	return run(args[0], args[1:]...)
}

// Run executes a command with privilege escalation unconditionally.
func Run(args ...string) error {
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
