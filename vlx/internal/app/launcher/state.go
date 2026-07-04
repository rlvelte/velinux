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
	Usage map[string]UsageEntry `json:"usage"`
}

type UsageEntry struct {
	Count    int   `json:"count"`
	LastUsed int64 `json:"last_used"`
}

// RankedEntry wraps Entry with usage info for sorting.
type RankedEntry struct {
	Entry
	Count    int
	LastUsed int64
}

func decodeState(_, _ string, data []byte) (*State, error) {
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	if s.Usage == nil {
		s.Usage = make(map[string]UsageEntry)
	}
	return &s, nil
}

// SaveState atomically writes the state file.
func SaveState(s *State) error {
	dir := filepath.Dir(filepath.Join(fsys.DataPath("vlx", "launcher"), "state.json"))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	tmp := filepath.Join(fsys.DataPath("vlx", "launcher"), "state.json") + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmp, filepath.Join(fsys.DataPath("vlx", "launcher"), "state.json"))
}

// Bump increments the usage count for an application.
func Bump(s *State, id string) {
	u := s.Usage[id]
	u.Count++
	u.LastUsed = time.Now().Unix()
	s.Usage[id] = u
}

// Rank returns entries sorted by the state rules.
func Rank(entries []Entry, state *State) []Entry {
	var ranked []RankedEntry
	for _, e := range entries {
		u := state.Usage[e.ID]
		ranked = append(ranked, RankedEntry{Entry: e, Count: u.Count, LastUsed: u.LastUsed})
	}

	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].Count != ranked[j].Count {
			return ranked[i].Count > ranked[j].Count
		}
		return ranked[i].Name < ranked[j].Name
	})

	var result []Entry
	for i := 0; i < len(ranked) && i < 10; i++ {
		result = append(result, ranked[i].Entry)
	}

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
