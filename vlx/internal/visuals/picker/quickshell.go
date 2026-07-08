package picker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/rlvelte/velinux/vlx/internal/core/guard"
)

// QuickshellPicker is a backend that uses quickshell picker.
type QuickshellPicker struct{}

// Available reports whether quickshell is installed.
func (q *QuickshellPicker) Available() bool {
	return guard.Binaries("quickshell") == nil
}

// Select prompts the user to choose one item via quickshell.
func (q *QuickshellPicker) Select(ctx context.Context, items []Item) (Item, error) {
	result, err := q.pick(ctx, items, "singlepicker")
	if err != nil {
		return Item{}, err
	}

	if len(result) == 0 {
		return Item{}, fmt.Errorf("picker cancelled")
	}

	return result[0], nil
}

// SelectMulti prompts the user to choose multiple items via quickshell multi.
func (q *QuickshellPicker) SelectMulti(ctx context.Context, items []Item) ([]Item, error) {
	return q.pick(ctx, items, "multipicker")
}

// SelectTwoStage prompts the user to choose an item and then a subitem via quickshell.
func (q *QuickshellPicker) SelectTwoStage(ctx context.Context, items []Item) (Item, error) {
	ts := fmt.Sprintf("%d", time.Now().UnixNano())
	itemsFile := filepath.Join("/dev/shm", "vlx-picker-"+ts+"-items")
	resultFile := filepath.Join("/dev/shm", "vlx-picker-"+ts+"-result")

	data, err := json.Marshal(items)
	if err != nil {
		return Item{}, fmt.Errorf("marshalling picker items: %w", err)
	}

	if err := os.WriteFile(itemsFile, data, 0644); err != nil {
		return Item{}, fmt.Errorf("writing picker items: %w", err)
	}
	defer os.Remove(itemsFile)

	if err := exec.CommandContext(ctx, "quickshell", "ipc", "call", "twostagepicker", "vlxOpen", itemsFile, resultFile).Run(); err != nil {
		return Item{}, fmt.Errorf("picker twostage: %w", err)
	}

	if err := q.waitForResult(ctx, resultFile); err != nil {
		return Item{}, err
	}
	defer os.Remove(resultFile)

	raw, err := os.ReadFile(resultFile)
	if err != nil {
		return Item{}, fmt.Errorf("reading picker result: %w", err)
	}

	var result struct {
		Item       Item  `json:"item"`
		Subcommand *Item `json:"subcommand"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return Item{}, fmt.Errorf("parsing picker result: %w", err)
	}

	if result.Subcommand != nil {
		return *result.Subcommand, nil
	}
	return result.Item, nil
}

func (q *QuickshellPicker) pick(ctx context.Context, items []Item, target string) ([]Item, error) {
	ts := fmt.Sprintf("%d", time.Now().UnixNano())
	itemsFile := filepath.Join("/dev/shm", "vlx-picker-"+ts+"-items")
	resultFile := filepath.Join("/dev/shm", "vlx-picker-"+ts+"-result")

	data, err := json.Marshal(items)
	if err != nil {
		return nil, fmt.Errorf("marshalling picker items: %w", err)
	}

	if err := os.WriteFile(itemsFile, data, 0644); err != nil {
		return nil, fmt.Errorf("writing picker items: %w", err)
	}
	defer os.Remove(itemsFile)

	if err := exec.CommandContext(ctx, "quickshell", "ipc", "call", target, "vlxOpen", itemsFile, resultFile).Run(); err != nil {
		return nil, fmt.Errorf("picker %s: %w", target, err)
	}

	if err := q.waitForResult(ctx, resultFile); err != nil {
		return nil, err
	}
	defer os.Remove(resultFile)

	raw, err := os.ReadFile(resultFile)
	if err != nil {
		return nil, fmt.Errorf("reading picker result: %w", err)
	}

	var result []Item
	if err := json.Unmarshal(raw, &result); err != nil {
		var single Item
		if err2 := json.Unmarshal(raw, &single); err2 != nil {
			return nil, fmt.Errorf("parsing picker result: %w", err)
		}

		result = append(result, single)
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
