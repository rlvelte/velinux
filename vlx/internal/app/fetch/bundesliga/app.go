package bundesliga

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/rlvelte/velinux/vlx/internal/core/guard"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
	"github.com/spf13/cobra"
)

const apiBase = "https://api.openligadb.de"

// setup validates all requirements for further processing.
func setup(_ *cobra.Command, _ []string) error {
	return guard.Network()
}

// Command returns the cobra command tree for vlx fetch Bundesliga.
func Command() *cobra.Command {
	root := &cobra.Command{
		Use:          "bundesliga",
		Short:        "Horribly bad Bundesliga tracker",
		Aliases:      []string{"buli", "bl"},
		SilenceUsage: !isatty.IsTerminal(os.Stdout.Fd()),
		PreRunE:      setup,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	tableCmd := &cobra.Command{
		Use:   "table",
		Short: "Show current table",
		RunE:  cmdTable,
	}

	tableCmd.Flags().String("league", "bl1", "league shortcut (bl1, bl2)")
	root.AddCommand(tableCmd)
	return root
}

func cmdTable(cmd *cobra.Command, _ []string) error {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)

	league, _ := cmd.Flags().GetString("league")
	url := fmt.Sprintf("%s/getbltable/%s/%d", apiBase, league, season())
	body, err := get(url)
	if err != nil {
		return fmt.Errorf("fetching table: %w", err)
	}

	var entries []TableEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		return fmt.Errorf("decoding table: %w", err)
	}

	rows := make([]TableRow, 0, len(entries))
	for i, e := range entries {
		rows = append(rows, TableRow{
			Position:      i + 1,
			Team:          e.TeamName,
			Points:        e.Points,
			GoalsFor:      e.Goals,
			GoalsAgainst:  e.OpponentGoals,
			Wins:          e.Won,
			Draws:         e.Draw,
			Losses:        e.Lost,
			MatchesPlayed: e.Matches,
		})
	}

	headers := []string{"#", "Team", "Pts", "Goals", "W/D/L"}
	var stringRows [][]string
	for _, r := range rows {
		stringRows = append(stringRows, []string{
			fmt.Sprintf("%d", r.Position),
			r.Team,
			fmt.Sprintf("%d", r.Points),
			fmt.Sprintf("%d:%d", r.GoalsFor, r.GoalsAgainst),
			fmt.Sprintf("%d/%d/%d", r.Wins, r.Draws, r.Losses),
		})
	}

	p.Table(headers, stringRows)
	return nil
}

func season() int {
	now := time.Now()
	if now.Month() >= time.August {
		return now.Year()
	}

	return now.Year() - 1
}

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
