package launcher

import (
	"encoding/json"
)

// Config holds user preferences for the launcher.
type Config struct {
	Ignore []string `json:"ignore"`
}

// IsIgnored checks whether an application ID is in the ignore list.
func IsIgnored(c *Config, id string) bool {
	for _, ignored := range c.Ignore {
		if ignored == id {
			return true
		}
	}
	return false
}

// State tracks how often each application is launched.
type State struct {
	Usage map[string]UsageEntry `json:"usage"`
}

// UsageEntry records a single application's launch statistics.
type UsageEntry struct {
	Count    int   `json:"count"`
	LastUsed int64 `json:"last_used"`
}

func decodeConfig(_ string, _ string, data []byte) (*Config, error) {
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func decodeState(_ string, _ string, data []byte) (*State, error) {
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	if s.Usage == nil {
		s.Usage = make(map[string]UsageEntry)
	}
	return &s, nil
}
