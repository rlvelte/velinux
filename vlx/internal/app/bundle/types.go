package bundle

import (
	"encoding/json"
	"errors"
)

type Repo struct {
	Alias string `json:"alias"`
	URL   string `json:"url"`
}

type Scheme struct {
	Name     string
	Repos    []Repo   `json:"repos"`
	Zypper   []string `json:"zypper"`
	Flatpak  []string `json:"flatpak"`
	PreHook  []string `json:"pre"`
	PostHook []string `json:"post"`
}

func decodeBundle(name, _ string, data []byte) (Scheme, error) {
	if name == "bundle.schema" {
		return Scheme{}, errors.New("skip schema file")
	}

	var s Scheme
	if err := json.Unmarshal(data, &s); err != nil {
		return Scheme{}, err
	}

	s.Name = name
	return s, nil
}
