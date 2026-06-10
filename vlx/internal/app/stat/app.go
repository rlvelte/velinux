package stat

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rvelte/vlx/internal/app/stat/sources"
	"github.com/rvelte/vlx/internal/system/guard"
	"github.com/spf13/cobra"
)

var allSources = []sources.Source{
	sources.Zypper{},
	sources.Flatpak{},
}

func ensure(_ *cobra.Command, _ []string) error {
	return errors.Join(guard.OS(), guard.Connection())
}

func Command() *cobra.Command {
	root := &cobra.Command{
		Use:               "stat",
		Short:             "System status and update checks",
		PersistentPreRunE: ensure,
		RunE:              cmdCheck,
	}

	root.AddCommand(
		&cobra.Command{
			Use:   "check",
			Short: "Check all sources for available updates",
			RunE:  cmdCheck,
		},
		&cobra.Command{
			Use:   "waybar",
			Short: "Output update status as waybar JSON",
			RunE:  cmdWaybar,
		},
	)

	return root
}

func cmdCheck(_ *cobra.Command, _ []string) error {
	state := &UpdateState{
		LastCheck: time.Now(),
		Sources:   make(map[string]SourceState),
	}

	total := 0
	for _, src := range allSources {
		count, pkgs, err := src.CheckUpdates()
		if err != nil {
			return fmt.Errorf("%s: %w", src.Name(), err)
		}

		state.Sources[src.Name()] = SourceState{
			Count:    count,
			Packages: pkgs,
		}
		total += count

		fmt.Printf("%s: %d update(s)\n", src.Name(), count)
		for _, p := range pkgs {
			fmt.Printf("  %s\n", p)
		}
	}

	fmt.Printf("\ntotal: %d update(s)\n", total)

	if err := writeState(state); err != nil {
		return fmt.Errorf("persist state: %w", err)
	}

	return nil
}

type waybarOutput struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
	Class   string `json:"class"`
}

func cmdWaybar(_ *cobra.Command, _ []string) error {
	cached, err := readState()
	if err != nil {
		return fmt.Errorf("read state: %w", err)
	}

	if cached == nil {
		for _, src := range allSources {
			count, pkgs, err := src.CheckUpdates()
			if err != nil {
				return fmt.Errorf("%s: %w", src.Name(), err)
			}
			cached = &UpdateState{
				LastCheck: time.Now(),
				Sources:   make(map[string]SourceState),
			}
			cached.Sources[src.Name()] = SourceState{Count: count, Packages: pkgs}
		}
		_ = writeState(cached)
	}

	total := 0
	var lines []string
	for _, src := range allSources {
		ss, ok := cached.Sources[src.Name()]
		if !ok {
			continue
		}
		total += ss.Count
		if ss.Count > 0 {
			lines = append(lines, fmt.Sprintf("%s: %s", src.Name(), strings.Join(ss.Packages, ", ")))
		}
	}

	class := "idle"
	if total > 0 {
		class = "updates"
	}

	out := waybarOutput{
		Text:    fmt.Sprintf("%d", total),
		Tooltip: strings.Join(lines, "\n"),
		Class:   class,
	}

	data, err := json.Marshal(out)
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}
