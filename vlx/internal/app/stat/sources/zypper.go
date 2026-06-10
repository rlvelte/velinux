package sources

import (
	"context"
	"encoding/xml"

	"github.com/rvelte/vlx/internal/core/cmd"
)

type Zypper struct{}

func (Zypper) Name() string { return "zypper" }

type zypperStream struct {
	XMLName      xml.Name           `xml:"stream"`
	UpdateStatus zypperUpdateStatus `xml:"update-status"`
}

type zypperUpdateStatus struct {
	Version string         `xml:"version,attr"`
	Updates []zypperUpdate `xml:"update-list>update"`
}

type zypperUpdate struct {
	Name       string `xml:"name,attr"`
	Edition    string `xml:"edition,attr"`
	EditionOld string `xml:"edition-old,attr"`
	Arch       string `xml:"arch,attr"`
	Kind       string `xml:"kind,attr"`
}

func (Zypper) CheckUpdates() (int, []string, error) {
	res, err := cmd.New("zypper", "--xmlout", "list-updates").RunCaptured(context.Background())
	if err != nil {
		return 0, nil, err
	}
	out := res.Text()

	var s zypperStream
	if err := xml.Unmarshal([]byte(out), &s); err != nil {
		return 0, nil, err
	}

	var pkgs []string
	for _, u := range s.UpdateStatus.Updates {
		pkgs = append(pkgs, u.Name)
	}

	return len(pkgs), pkgs, nil
}
