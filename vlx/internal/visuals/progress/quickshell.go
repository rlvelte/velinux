package progress

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	"github.com/rlvelte/velinux/vlx/internal/core/guard"
)

// state is the JSON payload shared between Go and the QML popup.
type state struct {
	Label      string  `json:"label"`
	Progress   float64 `json:"progress"`
	Indefinite bool    `json:"indefinite"`
	Done       bool    `json:"done"`
}

// QuickshellProgress is a backend that shows progress via a quickshell IPC popup.
type QuickshellProgress struct {
	dir     string
	label   string
	total   int
	current int
	active  bool
}

func (q *QuickshellProgress) Available() bool {
	return guard.Binaries("quickshell") == nil
}

// Start opens a quickshell progress popup. Total = 0 for indefinite spinner.
func (q *QuickshellProgress) Start(label string, total int) {
	ts := fmt.Sprintf("%d", time.Now().UnixNano())
	q.dir = filepath.Join("/dev/shm", "vlx-progress-"+ts)
	q.label = label
	q.total = total
	q.current = 0
	q.active = true

	q.cleanupStale()
	if err := os.MkdirAll(q.dir, 0755); err != nil {
		q.active = false
		return
	}

	q.writeState()
	_ = exec.CommandContext(context.Background(), "quickshell", "ipc", "call", "progress", "vlxOpen", q.dir).Run()
}

func (q *QuickshellProgress) cleanupStale() {
	entries, err := os.ReadDir("/dev/shm")
	if err != nil {
		return
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "vlx-progress-") {
			os.RemoveAll(filepath.Join("/dev/shm", e.Name()))
		}
	}
}

// Advance pulses the spinner (indefinite) or advances the counter (determinate).
func (q *QuickshellProgress) Advance(int) {
	if !q.active {
		return
	}

	q.writeState()
}

// SetProgress sets exact percentage 0.0–1.0.
func (q *QuickshellProgress) SetProgress(pct float64) {
	if !q.active {
		return
	}

	q.current = int(pct * float64(q.total))
	q.writeState()
}

// SetLabel updates the label message.
func (q *QuickshellProgress) SetLabel(label string) {
	if !q.active {
		return
	}

	q.label = label
	q.writeState()
}

// Stop writes a done state and cleans up the temp directory.
func (q *QuickshellProgress) Stop() {
	if !q.active {
		return
	}

	s := q.buildState()
	s.Done = true
	q.writeStateLocked(s)

	time.Sleep(250 * time.Millisecond) // let the QML poller read the Done state
	os.RemoveAll(q.dir)

	q.active = false
}

// writeState serialises the current progress state to disk.
func (q *QuickshellProgress) writeState() {
	q.writeStateLocked(q.buildState())
}

func (q *QuickshellProgress) buildState() state {
	var pct float64
	if q.total > 0 {
		pct = float64(q.current) / float64(q.total)
		if pct > 1.0 {
			pct = 1.0
		}
	}
	return state{
		Label:      q.label,
		Progress:   pct,
		Indefinite: q.total <= 0,
	}
}

func (q *QuickshellProgress) writeStateLocked(s state) {
	data, err := json.Marshal(s)
	if err != nil {
		return
	}

	_ = fsys.AtomicWrite(filepath.Join(q.dir, "state.json"), data)
}
