package launcher

import (
	"encoding/json"
)

type Config struct {
	Ignore []string `json:"ignore"`
}

func decodeConfig(_, _ string, data []byte) (*Config, error) {
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
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
