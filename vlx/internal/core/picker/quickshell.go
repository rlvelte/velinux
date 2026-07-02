package picker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var _ Variant = (*QuickshellPicker)(nil)

type QuickshellPicker struct{}

func (q *QuickshellPicker) Available() bool {
	_, err := exec.LookPath("quickshell")
	if err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	return exec.CommandContext(ctx, "quickshell", "ipc", "show").Run() == nil
}

func (q *QuickshellPicker) Select(ctx context.Context, items []string) (string, error) {
	result, err := q.pick(ctx, items, false)
	if err != nil {
		return "", err
	}
	if len(result) == 0 {
		return "", fmt.Errorf("picker cancelled")
	}
	return result[0], nil
}

func (q *QuickshellPicker) SelectMulti(ctx context.Context, items []string) ([]string, error) {
	return q.pick(ctx, items, true)
}

func (q *QuickshellPicker) pick(ctx context.Context, items []string, multi bool) ([]string, error) {
	ts := fmt.Sprintf("%d", time.Now().UnixNano())
	itemsFile := filepath.Join(os.TempDir(), "vlx-picker-"+ts+"-items")
	resultFile := filepath.Join(os.TempDir(), "vlx-picker-"+ts+"-result")

	if err := os.WriteFile(itemsFile, []byte(strings.Join(items, "\n")), 0644); err != nil {
		return nil, fmt.Errorf("writing picker items: %w", err)
	}
	defer os.Remove(itemsFile)

	cmd := "open"
	if multi {
		cmd = "openMulti"
	}
	if err := exec.CommandContext(ctx, "quickshell", "ipc", "call", "picker", cmd, itemsFile, resultFile).Run(); err != nil {
		return nil, fmt.Errorf("picker %s: %w", cmd, err)
	}

	if err := q.waitForResult(ctx, resultFile); err != nil {
		return nil, err
	}
	defer os.Remove(resultFile)

	data, err := os.ReadFile(resultFile)
	if err != nil {
		return nil, fmt.Errorf("reading picker result: %w", err)
	}

	var result []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}
	return result, nil
}

func (q *QuickshellPicker) waitForResult(ctx context.Context, resultFile string) error {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	timeout := time.After(120 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			return fmt.Errorf("picker timed out")
		case <-ticker.C:
			if _, err := os.Stat(resultFile); err == nil {
				return nil
			}
		}
	}
}
