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
func (e *Picker) Select(ctx context.Context, items []string) (string, error) {
	return e.variant.Select(ctx, items)
}

// SelectMulti prompts the user to choose multiple items.
func (e *Picker) SelectMulti(ctx context.Context, items []string) ([]string, error) {
	return e.variant.SelectMulti(ctx, items)
}

// ForceRofi forces the rofi backend.
func (e *Picker) ForceRofi() *Picker {
	e.variant = &RofiPicker{}
	return e
}

// ForceFzf forces the fzf backend.
func (e *Picker) ForceFzf() *Picker {
	e.variant = &FzfPicker{}
	return e
}

func auto() Variant {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		f := &FzfPicker{}
		if f.Available() {
			return f
		}
	}

	if os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != "" {
		r := &RofiPicker{}
		if r.Available() {
			return r
		}
	}

	return nil
}
