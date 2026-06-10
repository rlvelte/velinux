package printer

import (
	"github.com/rvelte/vlx/internal/core/printer/backends"
)

// Backend handles terminal out/inputs
type Backend interface {
	Available() bool                             // Available reports whether this backend can be used.
	Info(msg string)                             // Info prints an info message
	Success(msg string)                          // Success prints a success message
	Warn(msg string)                             // Warn prints a warning message
	Error(msg string)                            // Error prints a error message
	Header(msg string)                           // Header signals a new section
	Table(headers []string, rows [][]string)     // Table prints data in a tabular format.
	Confirm(msg string) bool                     // Confirm shows a simple confirmation dialog.
	Spinner(label string, fn func() error) error // Spinner shows a progress indicator.
}

// Printer is the unified printing engine.
type Printer struct {
	backend Backend // The selected Backend for this printer.
}

// New creates an engine with an auto-detected backend.
func New() *Printer {
	return &Printer{backend: detect()}
}

// Info prints an info message.
func (p *Printer) Info(msg string) {
	p.backend.Info(msg)
}

// Success prints a success message.
func (p *Printer) Success(msg string) {
	p.backend.Success(msg)
}

// Warn prints a warning message.
func (p *Printer) Warn(msg string) {
	p.backend.Warn(msg)
}

// Error prints an error message.
func (p *Printer) Error(msg string) {
	p.backend.Error(msg)
}

// Header signals a new section.
func (p *Printer) Header(msg string) {
	p.backend.Header(msg)
}

// Table prints data in a tabular format.
func (p *Printer) Table(headers []string, rows [][]string) {
	p.backend.Table(headers, rows)
}

// Confirm shows a simple confirmation dialog.
func (p *Printer) Confirm(msg string) bool {
	return p.backend.Confirm(msg)
}

// Spinner shows a progress indicator.
func (p *Printer) Spinner(label string, fn func() error) error {
	return p.backend.Spinner(label, fn)
}

// UseGum forces the gum backend.
func (p *Printer) UseGum() *Printer {
	p.backend = &backends.Gum{}
	return p
}

// UseBasic forces the basic backend.
func (p *Printer) UseBasic() *Printer {
	p.backend = &backends.Basic{}
	return p
}

func detect() Backend {
	g := &backends.Gum{}
	if g.Available() {
		return g
	}

	return &backends.Basic{}
}
