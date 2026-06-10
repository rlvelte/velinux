package picker

import (
	"context"
	"os"

	"github.com/rvelte/vlx/internal/core/picker/backends"
	"golang.org/x/term"
)

// Backend handles interactive item selection.
type Backend interface {
	Available() bool                                                   // Available reports whether this backend can be used.
	Select(ctx context.Context, items []string) (string, error)        // Select prompts the user to choose one item.
	SelectMulti(ctx context.Context, items []string) ([]string, error) // SelectMulti prompts the user to choose multiple items.
}

// Picker is the unified picking engine.
type Picker struct {
	backend Backend // The selected Backend for this picker.
}

// New creates an engine with an auto-detected backend.
func New() *Picker {
	return &Picker{
		backend: detect(),
	}
}

// Select prompts the user to choose one item.
func (e *Picker) Select(ctx context.Context, items []string) (string, error) {
	return e.backend.Select(ctx, items)
}

// SelectMulti prompts the user to choose multiple items.
func (e *Picker) SelectMulti(ctx context.Context, items []string) ([]string, error) {
	return e.backend.SelectMulti(ctx, items)
}

// UseRofi forces the rofi backend.
func (e *Picker) UseRofi() *Picker {
	e.backend = &backends.Rofi{}
	return e
}

// UseFzf forces the fzf backend.
func (e *Picker) UseFzf() *Picker {
	e.backend = &backends.Fzf{}
	return e
}

func detect() Backend {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		f := &backends.Fzf{}
		if f.Available() {
			return f
		}
	}

	if os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != "" {
		r := &backends.Rofi{}
		if r.Available() {
			return r
		}
	}

	return nil
}
