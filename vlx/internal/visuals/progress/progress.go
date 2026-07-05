package progress

import (
	"os"

	"github.com/mattn/go-isatty"
)

const ContextKey = "progress"

// Variant handles progress output methods.
type Variant interface {
	Available() bool
	Stream(message string)
}

// Progress is the unified progress engine.
type Progress struct {
	variant Variant
}

// New creates an engine with an auto-detected backend.
func New() *Progress {
	return &Progress{variant: auto()}
}

// Stream messages to the output source
func (p *Progress) Stream(message string) {
	p.variant.Stream(message)
}

// ForceFmt forces the basic backend.
func (p *Progress) ForceFmt() *Progress {
	p.variant = &FmtProgress{}
	return p
}

// ForceQuickshell forces the quickshell IPC backend.
func (p *Progress) ForceQuickshell() *Progress {
	p.variant = &QuickshellProgress{}
	return p
}

func auto() Variant {
	f := &FmtProgress{}
	if f.Available() && isatty.IsTerminal(os.Stdout.Fd()) {
		return f
	}

	q := &QuickshellProgress{}
	if q.Available() {
		return q
	}

	return nil
}
