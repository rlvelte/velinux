package printer

const ContextKey = "printer"

// Variant handles terminal out/inputs
type Variant interface {
	Available() bool                             // Available reports whether this backend can be used.
	Print(msg string)                            // Print prints a simple message
	Warn(msg string)                             // Warn prints a warning message
	Error(msg string)                            // Error prints a error message
	Table(headers []string, rows [][]string)     // Table prints data in a tabular format.
	Confirm(msg string, defaultYes bool) bool    // Confirm shows a simple confirmation dialog.
	Spinner(label string, fn func() error) error // Spinner shows a progress indicator.
}

// Printer is the unified printing engine.
type Printer struct {
	variant Variant // The selected Backend for these variants.
}

// New creates an engine with an auto-detected backend.
func New() *Printer {
	return &Printer{variant: auto()}
}

// Info prints an info message.
func (p *Printer) Info(msg string) {
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

// Confirm shows a simple confirmation dialog.
func (p *Printer) Confirm(msg string, defaultYes bool) bool {
	return p.variant.Confirm(msg, defaultYes)
}

// Spinner shows a progress indicator.
func (p *Printer) Spinner(label string, fn func() error) error {
	return p.variant.Spinner(label, fn)
}

// ForceFmt forces the basic backend.
func (p *Printer) ForceFmt() *Printer {
	p.variant = &FmtPrinter{}
	return p
}

func auto() Variant {
	return &FmtPrinter{}
}
