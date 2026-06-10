package backends

import (
	"context"
	"os/exec"
	"strings"
)

// Fzf is a backend that uses fzf command line picker.
type Fzf struct{}

// Available reports whether fzf is installed.
func (f *Fzf) Available() bool {
	_, err := exec.LookPath("fzf")
	return err == nil
}

// Select prompts the user to choose one item via fzf.
func (f *Fzf) Select(ctx context.Context, items []string) (string, error) {
	args := []string{"--prompt", "> ", "--height", "40%", "--reverse"}

	cmd := exec.CommandContext(ctx, "fzf", args...)
	cmd.Stdin = strings.NewReader(strings.Join(items, "\n"))

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// SelectMulti prompts the user to choose multiple items via fzf --multi.
func (f *Fzf) SelectMulti(ctx context.Context, items []string) ([]string, error) {
	args := []string{"--prompt", "> ", "--multi", "--height", "40%", "--reverse"}

	cmd := exec.CommandContext(ctx, "fzf", args...)
	cmd.Stdin = strings.NewReader(strings.Join(items, "\n"))

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var result []string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}
	return result, nil
}
