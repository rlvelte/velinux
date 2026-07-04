package launcher

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
)

// Config lives at ~/.config/vlx/launcher/config.json.
type Config struct {
	Ignore []string `json:"ignore"`
}

var configPath = filepath.Join(fsys.ConfigPath("vlx", "launcher"), "config.json")

// LoadConfig reads the ignore list, returning an empty one if absent.
func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}

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
