package pass

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Run executes a command with privilege escalation.
func Run(args ...string) error {
	if err := run("sudo", append([]string{"-n"}, args...)...); err == nil {
		return nil
	}

	if askpassPath := siblingBinary("vlxpass"); askpassPath != "" {
		sudoArgs := append([]string{"-A", "--preserve-env=WAYLAND_DISPLAY,HOME,XDG_RUNTIME_DIR"}, args...)
		cmd := exec.Command("sudo", sudoArgs...)
		cmd.Env = append(os.Environ(), "SUDO_ASKPASS="+askpassPath)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func siblingBinary(name string) string {
	exe, err := os.Executable()
	if err == nil {
		path := filepath.Join(filepath.Dir(exe), name)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	if p, err := exec.LookPath(name); err == nil {
		return p
	}

	return ""
}
