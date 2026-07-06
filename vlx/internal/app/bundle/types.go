package bundle

import (
	"encoding/json"
	"errors"
)

// Bundle is a bundle file.
type Bundle struct {
	path     string
	filename string
	Info     Info     `json:"info"`
	Repos    []Repo   `json:"repos"`
	Zypper   []string `json:"zypper"`
	Flatpak  []string `json:"flatpak"`
	PreHook  []string `json:"pre"`
	PostHook []string `json:"post"`
}

// Repo is a repository.
type Repo struct {
	Alias string `json:"alias"`
	URL   string `json:"url"`
}

// Info is a bundle info.
type Info struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
	Icon        string `json:"icon"`
}

// decodeBundle decodes a bundle file.
func decodeBundle(name, path string, data []byte) (Bundle, error) {
	if name == "bundle.schema" {
		return Bundle{}, errors.New("skip schema file")
	}

	var s Bundle
	if err := json.Unmarshal(data, &s); err != nil {
		return Bundle{}, err
	}

	s.path = path
	s.filename = name
	return s, nil
}
