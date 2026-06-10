package data

import (
	"strings"
	"time"
)

// CmdResult holds the output and metadata of a completed command.
type CmdResult struct {
	Stdout   []byte
	Stderr   []byte
	ExitCode int
	Duration time.Duration
}

// Text returns the stdout as a trimmed string.
func (r *CmdResult) Text() string {
	if r == nil {
		return ""
	}
	return strings.TrimSpace(string(r.Stdout))
}

// Lines returns the stdout split into non-empty lines.
func (r *CmdResult) Lines() []string {
	if r == nil || len(r.Stdout) == 0 {
		return nil
	}

	raw := strings.Split(string(r.Stdout), "\n")
	var lines []string
	for _, l := range raw {
		l = strings.TrimSpace(l)
		if l != "" {
			lines = append(lines, l)
		}
	}

	return lines
}

// Success reports whether the command exited with code 0.
func (r *CmdResult) Success() bool {
	return r != nil && r.ExitCode == 0
}
