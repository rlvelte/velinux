package printer

import (
	"encoding/json"
	"os"

	"github.com/mattn/go-isatty"
)

const ContextKey = "printer"

// Variant handles terminal out/inputs
type Variant interface {
	Success(msg string)
	Error(err error)
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

// Success prints an info message.
func (p *Printer) Success(msg string) {
	p.variant.Success(msg)
}

// Error prints an error message.
func (p *Printer) Error(err error) {
	p.variant.Error(err)
}

// Table prints data in a tabular format.
func (p *Printer) Table(headers []string, rows [][]string) {
	p.variant.Table(headers, rows)
}

// ForceFmt forces the basic backend.
func (p *Printer) ForceFmt() {
	p.variant = &FmtPrinter{}
}

// ForceJSON forces the JSON backend.
func (p *Printer) ForceJSON() {
	p.variant = &JsonPrinter{encoder: json.NewEncoder(os.Stdout)}
}

func auto() Variant {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		return &FmtPrinter{}
	}

	return &JsonPrinter{encoder: json.NewEncoder(os.Stdout)}
}
