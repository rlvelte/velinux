package picker

import (
	"context"
	"os"

	"github.com/mattn/go-isatty"
)

const ContextKey = "picker"

// Item represents a selectable item in the picker.
type Item struct {
	Icon        string `json:"icon"`
	Header      string `json:"header"`
	Description string `json:"description"`
}

// Variant handles interactive item selection.
type Variant interface {
	Available() bool                                               // Available reports whether this backend can be used.
	Select(ctx context.Context, items []Item) (Item, error)        // Select prompts the user to choose one item.
	SelectMulti(ctx context.Context, items []Item) ([]Item, error) // SelectMulti prompts the user to choose multiple items.
}

// Picker is the unified picking engine.
type Picker struct {
	variant Variant // The selected Variant for this picker.
}

// New creates an engine with an auto-detected backend.
// Returns nil if no backend is available.
func New() *Picker {
	v := auto()
	if v == nil {
		return nil
	}

	return &Picker{variant: v}
}

// Select prompts the user to choose one item.
func (p *Picker) Select(ctx context.Context, items []Item) (Item, error) {
	return p.variant.Select(ctx, items)
}

// SelectMulti prompts the user to choose multiple items.
func (p *Picker) SelectMulti(ctx context.Context, items []Item) ([]Item, error) {
	return p.variant.SelectMulti(ctx, items)
}

// ForceQuickshell forces the quickshell IPC backend.
func (p *Picker) ForceQuickshell() *Picker {
	p.variant = &QuickshellPicker{}
	return p
}

// ForceFzf forces the fzf backend.
func (p *Picker) ForceFzf() *Picker {
	p.variant = &FzfPicker{}
	return p
}

func auto() Variant {
	f := &FzfPicker{}
	if f.Available() && isatty.IsTerminal(os.Stdout.Fd()) {
		return f
	}

	q := &QuickshellPicker{}
	if q.Available() {
		return q
	}

	return nil
}
