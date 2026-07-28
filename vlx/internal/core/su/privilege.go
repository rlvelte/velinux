package su

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"unicode"

	"github.com/rlvelte/velinux/vlx/internal/core/logs"
)

var (
	reQuoted    = regexp.MustCompile(`'[^']*'|"[^"]*"`)
	reDashDash  = regexp.MustCompile(`\s--\s`)
	reOperators = regexp.MustCompile(`&&|\|\||[|;\n]`)
	reEnvPrefix = regexp.MustCompile(`^\w+=.*?\s+`)
)

// ShouldEscalate reports whether a shell command string likely needs privilege escalation on the host.
func ShouldEscalate(cmd string) bool {
	s := reQuoted.ReplaceAllString(cmd, "")
	s = removeAfterDashDash(s)

	for _, segment := range reOperators.Split(s, -1) {
		segment = strings.TrimSpace(segment)
		segment = stripEnvPrefix(segment)

		lower := strings.ToLower(segment)
		if strings.HasPrefix(lower, "sudo ") || strings.HasPrefix(lower, "su ") {
			return true
		}
	}

	return false
}

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

func removeAfterDashDash(s string) string {
	locs := reDashDash.FindAllStringIndex(s, -1)
	if len(locs) == 0 {
		return s
	}

	last := locs[len(locs)-1]
	return s[:last[0]]
}

func stripEnvPrefix(s string) string {
	for {
		loc := reEnvPrefix.FindStringIndex(s)
		if loc == nil || loc[0] != 0 {
			break
		}

		eq := strings.IndexByte(s[:loc[1]], '=')
		if eq < 0 || eq == 0 {
			break
		}

		name := s[:eq]
		if !isIdent(name) {
			break
		}

		s = s[loc[1]:]
	}

	return s
}

func isIdent(s string) bool {
	if len(s) == 0 {
		return false
	}

	for i, r := range s {
		if i == 0 && unicode.IsDigit(r) {
			return false
		}

		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') {
			return false
		}
	}

	return true
}
