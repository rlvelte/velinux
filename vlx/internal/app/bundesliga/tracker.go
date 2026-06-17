package bundesliga

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	vlxhttp "github.com/rlvelte/velinux/vlx/internal/core/http"
	"github.com/rlvelte/velinux/vlx/internal/core/notify"
)

type matchState struct {
	MatchID   int `json:"matchID"`
	HomeScore int `json:"homeScore"`
	AwayScore int `json:"awayScore"`
}

func poll(http *vlxhttp.HTTP) error {
	store := fsys.NewStore(fsys.ConfigPath("vlx", "bundesliga"), decodeConfig, ".json")
	cfg, err := store.Get("config")
	if err != nil || cfg.Team.TeamID == 0 {
		return nil
	}

	state := loadState()

	for _, league := range []string{"bl1", "bl2"} {
		group, err := currentGroup(http, league)
		if err != nil {
			continue
		}

		url := fmt.Sprintf("%s/getmatchdata/%s/%d/%d", apiBase, league, season(), group)
		body, err := http.Get(url)
		if err != nil {
			continue
		}

		var matches []Match
		if err := json.Unmarshal(body, &matches); err != nil {
			continue
		}

		for _, m := range matches {
			if m.Team1.TeamID != cfg.Team.TeamID && m.Team2.TeamID != cfg.Team.TeamID {
				continue
			}
			if m.MatchIsFinished {
				continue
			}
			if len(m.MatchResults) == 0 {
				return nil
			}

			home, away := currentScore(m)

			if state.MatchID == m.MatchID && state.HomeScore == home && state.AwayScore == away {
				return nil
			}

			if cfg.Notifications {
				if err := notify.New().Send(
					fmt.Sprintf("%s %d:%d %s", m.Team1.ShortName, home, away, m.Team2.ShortName),
					&notify.Details{
						Title:   "TOR!",
						Urgency: "normal",
						Timeout: 5000,
					},
				); err != nil {
					return fmt.Errorf("notification failed: %w", err)
				}
			}

			saveState(m.MatchID, home, away)
			return nil
		}
	}

	return nil
}

func currentGroup(http *vlxhttp.HTTP, league string) (int, error) {
	body, err := http.Get(fmt.Sprintf("%s/getcurrentgroup/%s", apiBase, league))
	if err != nil {
		return 0, err
	}
	var grp Group
	if err := json.Unmarshal(body, &grp); err != nil {
		return 0, err
	}
	return grp.GroupOrderID, nil
}

func currentScore(m Match) (int, int) {
	for _, r := range m.MatchResults {
		if r.ResultTypeID == 2 {
			return r.PointsTeam1, r.PointsTeam2
		}
	}
	return 0, 0
}

func statePath() string {
	return fsys.CachePath("vlx", "bundesliga", "state.json")
}

func loadState() matchState {
	var s matchState
	_ = os.MkdirAll(fsys.CachePath("vlx", "bundesliga"), 0755)
	data, err := os.ReadFile(statePath())
	if err != nil {
		return s
	}
	_ = json.Unmarshal(data, &s)
	return s
}

func saveState(matchID, home, away int) {
	_ = os.MkdirAll(fsys.CachePath("vlx", "bundesliga"), 0755)
	s := matchState{MatchID: matchID, HomeScore: home, AwayScore: away}
	data, _ := json.MarshalIndent(s, "", "  ")
	_ = os.WriteFile(statePath(), data, 0644)
}
