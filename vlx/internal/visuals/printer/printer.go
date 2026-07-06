package printer

import (
	"encoding/json"
	"os"

	"github.com/mattn/go-isatty"
)

const ContextKey = "printer"

// Variant handles terminal out/inputs
type Variant interface {
	Print(msg string)
	Warn(msg string)
	Error(msg string)
	Table(headers []string, rows [][]string)
}

// Printer is the unified printing engine.
type Printer struct {
	variant Variant // The selected Backend for these variants.
}

// New creates an engine with an auto-detected backend.
func New() *Printer {
	return &Printer{variant: auto()}
}

// Print prints an info message.
func (p *Printer) Print(msg string) {
	p.variant.Print(msg)
}

// Warn prints a warning message.
func (p *Printer) Warn(msg string) {
	p.variant.Warn(msg)
}

// Error prints an error message.
func (p *Printer) Error(msg string) {
	p.variant.Error(msg)
}

// Table prints data in a tabular format.
func (p *Printer) Table(headers []string, rows [][]string) {
	p.variant.Table(headers, rows)
}

// ForceFmt forces the basic backend.
func (p *Printer) ForceFmt() *Printer {
	p.variant = &FmtPrinter{}
	return p
}

// ForceJSON forces the JSON backend.
func (p *Printer) ForceJSON() *Printer {
	p.variant = &JSONPrinter{encoder: json.NewEncoder(os.Stdout)}
	return p
}

func auto() Variant {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		return &FmtPrinter{}
	}

	return &JSONPrinter{encoder: json.NewEncoder(os.Stdout)}
}
