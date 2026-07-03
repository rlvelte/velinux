package launcher

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
)

// State tracks how often each application is launched.
type State struct {
	Usage map[string]UsageEntry `json:"usage"` // key = desktop entry ID
}

type UsageEntry struct {
	Count    int   `json:"count"`
	LastUsed int64 `json:"last_used"` // Unix timestamp
}

var statePath = filepath.Join(fsys.DataPath("vlx", "launcher"), "state.json")

// LoadState reads the state file, creating a fresh one if absent.
func LoadState() (*State, error) {
	data, err := os.ReadFile(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{Usage: make(map[string]UsageEntry)}, nil
		}
		return nil, err
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return &State{Usage: make(map[string]UsageEntry)}, nil
	}
	if s.Usage == nil {
		s.Usage = make(map[string]UsageEntry)
	}
	return &s, nil
}

// SaveState atomically writes the state file.
func SaveState(s *State) error {
	dir := filepath.Dir(statePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	tmp := statePath + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, statePath)
}

// Bump increments the usage count for an application.
func Bump(s *State, id string) {
	u := s.Usage[id]
	u.Count++
	u.LastUsed = time.Now().Unix()
	s.Usage[id] = u
}

// RankedEntry wraps Entry with usage info for sorting.
type RankedEntry struct {
	Entry
	Count    int
	LastUsed int64
}

// Rank returns entries sorted by the state rules.
func Rank(entries []Entry, state *State) []Entry {
	var ranked []RankedEntry
	for _, e := range entries {
		u := state.Usage[e.ID]
		ranked = append(ranked, RankedEntry{Entry: e, Count: u.Count, LastUsed: u.LastUsed})
	}

	// Sort: count desc, name asc
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].Count != ranked[j].Count {
			return ranked[i].Count > ranked[j].Count
		}
		return ranked[i].Name < ranked[j].Name
	})

	var result []Entry
	// Take top 10
	for i := 0; i < len(ranked) && i < 10; i++ {
		result = append(result, ranked[i].Entry)
	}

	// Append remaining entries in ABC order
	if len(ranked) > 10 {
		rest := ranked[10:]
		sort.Slice(rest, func(i, j int) bool {
			return rest[i].Name < rest[j].Name
		})
		for _, e := range rest {
			result = append(result, e.Entry)
		}
	}

	return result
}
