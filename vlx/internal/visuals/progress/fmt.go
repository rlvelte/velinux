package progress

import (
	"fmt"
	"strings"
)

// spinner is a braille spinner sequence for indefinite progress.
const spinner = "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏"

// FmtProgress is a TUI backend that uses fmt to render progress spinners/bars.
type FmtProgress struct {
	label   string
	total   int
	current int
	pos     int
	started bool
}

func (f *FmtProgress) Available() bool {
	return true
}

// Start begins a progress display. Total = 0 for spinner, > 0 for bar.
func (f *FmtProgress) Start(label string, total int) {
	f.label = label
	f.total = total
	f.current = 0
	f.pos = 0
	f.started = true

	if total > 0 {
		f.renderBar()
	} else {
		f.renderSpinner()
	}
}

// Advance pulses the spinner (indefinite) or has no direct effect (determinate is step-based).
func (f *FmtProgress) Advance(int) {
	if !f.started {
		return
	}

	f.pos = (f.pos + 1) % len(spinner)
	f.renderSpinner()
}

// SetProgress updates the bar to an exact percentage 0.0–1.0.
func (f *FmtProgress) SetProgress(pct float64) {
	if !f.started {
		return
	}

	f.current = int(pct * float64(f.total))
	f.renderBar()
}

// SetLabel updates the displayed label.
func (f *FmtProgress) SetLabel(label string) {
	if !f.started {
		return
	}

	f.label = label
	if f.total > 0 {
		f.renderBar()
	} else {
		f.renderSpinner()
	}
}

// Stop finishes the display with a newline and final state.
func (f *FmtProgress) Stop() {
	if !f.started {
		return
	}

	if f.total > 0 {
		f.SetProgress(1.0)
		fmt.Print("\033[999;0H\n")
	} else {
		fmt.Print("\033[999;0H\033[K")
		fmt.Println("✓ " + f.label)
	}

	f.started = false
}

func (f *FmtProgress) renderBar() {
	const barWidth = 40
	pct := float64(f.current) / float64(f.total)
	if pct > 1.0 {
		pct = 1.0
	}

	if pct < 0 {
		pct = 0
	}

	filled := int(pct * barWidth)
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	fmt.Print("\0337")
	fmt.Printf("\033[999;0H\033[K[%s] %3.0f%% %s", bar, pct*100.0, f.label)
	fmt.Print("\0338")
}

func (f *FmtProgress) renderSpinner() {
	ch := string(spinner[f.pos])
	fmt.Print("\0337")
	fmt.Printf("\033[999;0H\033[K%s %s", ch, f.label)
	fmt.Print("\0338")
}
