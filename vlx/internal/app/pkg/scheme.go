package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tailscale/hujson"
)

type Repo struct {
	Alias string `json:"alias"`
	URL   string `json:"url"`
}

type Scheme struct {
	Name     string   `json:"-"`
	Repos    []Repo   `json:"repos"`
	Zypper   []string `json:"zypper"`
	Flatpak  []string `json:"flatpak"`
	PreHook  []string `json:"pre"`
	PostHook []string `json:"post"`
}

func Load(name string) (*Scheme, error) {
	for _, ext := range []string{".json", ".jsonc"} {
		path := filepath.Join(stateDir(), name+ext)
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("read scheme %s: %w", name, err)
		}

		cleaned, err := hujson.Standardize(data)
		if err != nil {
			return nil, fmt.Errorf("parse scheme %s: %w", name, err)
		}

		var s Scheme
		if err := json.Unmarshal(cleaned, &s); err != nil {
			return nil, fmt.Errorf("parse scheme %s: %w", name, err)
		}
		s.Name = name
		return &s, nil
	}

	return nil, fmt.Errorf("scheme not found: %s", name)
}

func List() ([]*Scheme, error) {
	dir := stateDir()

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read schemes dir: %w", err)
	}

	var schemes []*Scheme
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".json") && !strings.HasSuffix(name, ".jsonc") {
			continue
		}
		name = strings.TrimSuffix(name, ".jsonc")
		name = strings.TrimSuffix(name, ".json")
		if strings.HasSuffix(name, ".schema") {
			continue
		}

		s, err := Load(name)
		if err != nil {
			continue
		}
		schemes = append(schemes, s)
	}

	return schemes, nil
}
