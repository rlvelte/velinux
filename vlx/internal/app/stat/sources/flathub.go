package sources

import (
	"context"
	"strings"

	"github.com/rvelte/vlx/internal/core/cmd"
)

type Flatpak struct{}

func (Flatpak) Name() string { return "flatpak" }

func (Flatpak) CheckUpdates() (int, []string, error) {
	res, err := cmd.New("flatpak", "remote-ls", "flathub", "--updates").RunCaptured(context.Background())
	if err != nil {
		return 0, nil, err
	}

	out := res.Text()
	lines := parseFlatpakUpdates(out)

	return len(lines), lines, nil
}

func parseFlatpakUpdates(output string) []string {
	var pkgs []string
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Ref") {
			continue
		}

		pkgs = append(pkgs, line)
	}

	return pkgs
}
