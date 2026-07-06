package picker

import (
	"context"
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
)

const ContextKey = "picker"

// Item represents a selectable item in the picker.
type Item struct {
	Icon        string `json:"icon"`
	Header      string `json:"header"`
	Description string `json:"description"`
	Subitems    []Item `json:"subitem"`
}

// Variant handles interactive item selection.
type Variant interface {
	Available() bool
	Select(ctx context.Context, items []Item) (Item, error)
	SelectMulti(ctx context.Context, items []Item) ([]Item, error)
	SelectTwoStage(ctx context.Context, items []Item) (Item, error)
}

// Picker is the unified picking engine.
type Picker struct {
	variant Variant
}

// New creates an engine with an auto-detected backend.
// Returns nil if no backend is available.
func New() (*Picker, error) {
	v, err := auto()
	if err != nil {
		return nil, err
	}

	return &Picker{variant: v}, nil
}

// Select prompts the user to choose one item.
func (p *Picker) Select(ctx context.Context, items []Item) (Item, error) {
	return p.variant.Select(ctx, items)
}

// SelectMulti prompts the user to choose multiple items via quickshell multi.
func (p *Picker) SelectMulti(ctx context.Context, items []Item) ([]Item, error) {
	return p.variant.SelectMulti(ctx, items)
}

// SelectTwoStage prompts the user to choose an item and then a subcommand.
func (p *Picker) SelectTwoStage(ctx context.Context, items []Item) (Item, error) {
	return p.variant.SelectTwoStage(ctx, items)
}

// ForceQuickshell forces the quickshell IPC backend.
func (p *Picker) ForceQuickshell() *Picker {
	p.variant = &Quickshell{}
	return p
}

// ForceFzf forces the fzf backend.
func (p *Picker) ForceFzf() *Picker {
	p.variant = &Fzf{}
	return p
}

func auto() (Variant, error) {
	f := &Fzf{}
	if f.Available() && isatty.IsTerminal(os.Stdout.Fd()) {
		return f, nil
	}

	q := &Quickshell{}
	if q.Available() {
		return q, nil
	}

	return nil, fmt.Errorf("no picker available")
}
