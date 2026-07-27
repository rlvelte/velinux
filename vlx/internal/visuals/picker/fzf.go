package picker

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rlvelte/velinux/vlx/internal/core/guard"
)

// FzfPicker is a backend that uses fzf command line picker.
type FzfPicker struct{}

// Available reports whether fzf is installed.
func (f *FzfPicker) Available() bool {
	return guard.Binaries("fzf") == nil
}

// Select prompts the user to choose one item via fzf.
func (f *FzfPicker) Select(ctx context.Context, items []Item) (Item, error) {
	selected, err := f.pick(ctx, items, false, false)
	if err != nil {
		return Item{}, err
	}

	if len(selected) == 0 {
		return Item{}, nil
	}

	return selected[0], nil
}

// SelectMulti prompts the user to choose multiple items via fzf --multi.
func (f *FzfPicker) SelectMulti(ctx context.Context, items []Item) ([]Item, error) {
	return f.pick(ctx, items, true, false)
}

// SelectTwoStage prompts the user to choose an item and then a subitem via fzf.
func (f *FzfPicker) SelectTwoStage(ctx context.Context, items []Item) (Item, error) {
	selected, err := f.pick(ctx, items, false, false)
	if err != nil {
		return Item{}, err
	}
	if len(selected) == 0 {
		return Item{}, nil
	}

	parent := selected[0]
	if len(parent.Subitems) == 0 {
		return parent, nil
	}

	sub, err := f.pick(ctx, parent.Subitems, false, false)
	if err != nil {
		return Item{}, err
	}
	if len(sub) == 0 {
		return parent, nil
	}

	return sub[0], nil
}

// SelectGrouped prompts the user to choose one item from a grouped list via fzf.
func (f *FzfPicker) SelectGrouped(ctx context.Context, items []Item) (Item, error) {
	selected, err := f.pick(ctx, items, false, true)
	if err != nil {
		return Item{}, err
	}
	if len(selected) == 0 {
		return Item{}, nil
	}
	return selected[0], nil
}

func (f *FzfPicker) pick(ctx context.Context, items []Item, multi bool, grouped bool) ([]Item, error) {
	lines := make([]string, len(items))
	for i, item := range items {
		s := item.Header
		if item.Description != "" {
			s += "  " + item.Description
		}
		if grouped && item.Group != "" {
			s = item.Group + ": " + s
		}

		lines[i] = fmt.Sprintf("%d\t%s", i, s)
	}

	args := []string{"--prompt", "> ", "--height", "40%", "--reverse", "--delimiter=\t", "--with-nth=2.."}
	if multi {
		args = append(args, "--multi")
	}

	cmd := exec.CommandContext(ctx, "fzf", args...)
	cmd.Stdin = strings.NewReader(strings.Join(lines, "\n"))

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var result []Item
	for line := range strings.SplitSeq(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "\t", 2)

		var idx int
		if _, err := fmt.Sscanf(parts[0], "%d", &idx); err == nil && idx >= 0 && idx < len(items) {
			result = append(result, items[idx])
		}
	}

	return result, nil
}
