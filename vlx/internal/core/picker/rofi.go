package picker

import (
	"context"
	"os/exec"
	"strings"
)

// RofiPicker is a backend that uses rofi graphical picker.
type RofiPicker struct {
	ConfigPath string
}

// Available reports whether rofi is installed.
func (r *RofiPicker) Available() bool {
	_, err := exec.LookPath("rofi")
	return err == nil
}

// Select prompts the user to choose one item via rofi.
func (r *RofiPicker) Select(ctx context.Context, items []string) (string, error) {
	args := r.baseArgs()
	args = append(args, "-dmenu")

	cmd := exec.CommandContext(ctx, "rofi", args...)
	cmd.Stdin = strings.NewReader(strings.Join(items, "\n"))

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// SelectMulti prompts the user to choose multiple items via rofi.
func (r *RofiPicker) SelectMulti(ctx context.Context, items []string) ([]string, error) {
	args := r.baseArgs()
	args = append(args, "-dmenu", "-multi-select")

	cmd := exec.CommandContext(ctx, "rofi", args...)
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

func (r *RofiPicker) baseArgs() []string {
	if r.ConfigPath != "" {
		return []string{"-config", r.ConfigPath}
	}
	return nil
}
