package picker

import (
	"context"
	"os"

	"golang.org/x/term"
)

const ContextKey = "picker"

// Variant handles interactive item selection.
type Variant interface {
	Available() bool                                                   // Available reports whether this backend can be used.
	Select(ctx context.Context, items []string) (string, error)        // Select prompts the user to choose one item.
	SelectMulti(ctx context.Context, items []string) ([]string, error) // SelectMulti prompts the user to choose multiple items.
}

// Picker is the unified picking engine.
type Picker struct {
	variant Variant // The selected Variant for this picker.
}

// New creates an engine with an auto-detected backend.
func New() *Picker {
	return &Picker{
		variant: auto(),
	}
}

// Select prompts the user to choose one item.
func (p *Picker) Select(ctx context.Context, items []string) (string, error) {
	return p.variant.Select(ctx, items)
}

// SelectMulti prompts the user to choose multiple items.
func (p *Picker) SelectMulti(ctx context.Context, items []string) ([]string, error) {
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
	if term.IsTerminal(int(os.Stdout.Fd())) {
		f := &FzfPicker{}
		if f.Available() {
			return f
		}
	}

	if os.Getenv("WAYLAND_DISPLAY") != "" || os.Getenv("DISPLAY") != "" {
		q := &QuickshellPicker{}
		if q.Available() {
			return q
		}
	}

	return nil
}
