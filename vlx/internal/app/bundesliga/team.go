package bundesliga

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	vlxhttp "github.com/rlvelte/velinux/vlx/internal/core/http"
	"github.com/rlvelte/velinux/vlx/internal/core/picker"
)

func selectTeam(ctx context.Context, http *vlxhttp.HTTP, cfg *Config) error {
	var allTeams []TeamInfo
	for _, league := range []string{"bl1", "bl2"} {
		url := fmt.Sprintf("%s/getavailableteams/%s/%d", apiBase, league, season())
		body, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("fetching %s teams: %w", league, err)
		}

		var teams []TeamInfo
		if err := json.Unmarshal(body, &teams); err != nil {
			return fmt.Errorf("decoding %s teams: %w", league, err)
		}
		allTeams = append(allTeams, teams...)
	}

	if len(allTeams) == 0 {
		return fmt.Errorf("no teams available")
	}

	names := make([]string, len(allTeams))
	for i, t := range allTeams {
		names[i] = fmt.Sprintf("%s (%s)", t.TeamName, t.ShortName)
	}

	sel := picker.New()
	chosen, err := sel.Select(ctx, names)
	if err != nil {
		return err
	}

	var selected TeamInfo
	for i, name := range names {
		if name == chosen {
			selected = allTeams[i]
			break
		}
	}

	if selected.TeamID == 0 {
		return fmt.Errorf("invalid selection")
	}

	cfg.Team = selected

	return nil
}

func saveConfig(cfg Config) error {
	path := fsys.ConfigPath("vlx", "bundesliga", "config.json")
	if err := os.MkdirAll(fsys.ConfigPath("vlx", "bundesliga"), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
