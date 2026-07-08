package progress

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
)

const ContextKey = "progress"

// Variant handles progress output display.
type Variant interface {
	Available() bool
	Start(label string, total int)
	Advance(n int)
	SetProgress(pct float64)
	SetLabel(label string)
	Stop()
}

// Progress is the unified progress engine.
type Progress struct {
	variant Variant
	current int
	total   int
}

// New creates an engine with an auto-detected backend.
func New() (*Progress, error) {
	v, err := auto()
	if err != nil {
		return nil, err
	}

	return &Progress{variant: v}, nil
}

// Start begins the progress display.
// Total = 0 for indefinite spinner, > 0 for determinate progress with that many steps.
func (p *Progress) Start(label string, total int) {
	p.current = 0
	p.total = total
	p.variant.Start(label, total)
}

// Advance advances the step count by n.
func (p *Progress) Advance(n int) {
	p.current += n
	if p.total > 0 {
		pct := float64(p.current) / float64(p.total)
		if pct > 1.0 {
			pct = 1.0
		}

		p.variant.SetProgress(pct)
	} else {
		p.variant.Advance(n)
	}
}

// SetProgress sets the progress to an exact percentage (0.0–1.0).
func (p *Progress) SetProgress(pct float64) {
	if pct < 0 {
		pct = 0
	}

	if pct > 1.0 {
		pct = 1.0
	}

	p.current = int(pct * float64(p.total))
	p.variant.SetProgress(pct)
}

// SetLabel updates the message displayed alongside the progress.
func (p *Progress) SetLabel(label string) {
	p.variant.SetLabel(label)
}

// Stop completes and cleans up the progress display.
func (p *Progress) Stop() {
	p.variant.Stop()
}

// ForceFmt forces the basic fmt backend.
func (p *Progress) ForceFmt() *Progress {
	p.variant = &FmtProgress{}
	return p
}

// ForceQuickshell forces the quickshell IPC backend.
func (p *Progress) ForceQuickshell() *Progress {
	p.variant = &QuickshellProgress{}
	return p
}

func auto() (Variant, error) {
	f := &FmtProgress{}
	if f.Available() && isatty.IsTerminal(os.Stdout.Fd()) {
		return f, nil
	}

	q := &QuickshellProgress{}
	if q.Available() {
		return q, nil
	}

	return nil, fmt.Errorf("progress: no backend available")
}
