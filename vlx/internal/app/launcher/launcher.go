package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	"github.com/rlvelte/velinux/vlx/internal/core/logs"
	"github.com/rlvelte/velinux/vlx/internal/visuals/notify"
	"github.com/rlvelte/velinux/vlx/internal/visuals/picker"
	"github.com/spf13/cobra"
)

// Command returns the launcher cobra command.
func Command() *cobra.Command {
	return &cobra.Command{
		Use:          "launcher",
		Short:        "Horribly bad application launcher",
		Long:         "Scan desktop entries and launch via the quickshell/fzf picker.",
		Aliases:      []string{"launch"},
		Args:         cobra.NoArgs,
		SilenceUsage: !isatty.IsTerminal(os.Stdout.Fd()),
		RunE:         cmdLaunch,
	}
}

// cmdLaunch launches the application picker.
func cmdLaunch(cmd *cobra.Command, _ []string) error {
	n, ok := cmd.Context().Value(notify.ContextKey).(*notify.Notify)
	if !ok {
		return fmt.Errorf("launcher: notify not configured")
	}

	cfg, err := fsys.GetJSON(fsys.ConfigPath("vlx", "launcher"), "config", decodeConfig)
	if err != nil {
		return fmt.Errorf("launcher config: %w", err)
	}

	state, err := fsys.GetJSON(fsys.DataPath("vlx", "launcher"), "state", decodeState)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("launcher state: %w", err)
		}

		state = &State{Usage: make(map[string]UsageEntry)}
		_ = fsys.SetJSON(fsys.DataPath("vlx", "launcher"), "state", state)
	}

	entries := Scan()
	var visible []Entry
	for _, e := range entries {
		if e.NoDisplay || e.Hidden || IsIgnored(cfg, e.ID) {
			continue
		}
		visible = append(visible, e)
	}

	ranked := rank(visible, state)

	items := make([]picker.Item, 0, len(ranked))
	entryOf := make(map[string]*Entry, len(ranked))
	for i := range ranked {
		e := &ranked[i]
		entryOf[e.Name] = e
		items = append(items, picker.Item{
			Icon:        parseIcon(e.Icon),
			Header:      e.Name,
			Description: e.Comment,
		})
	}

	pkr, err := picker.New()
	if err != nil {
		return err
	}

	selected, err := pkr.Select(cmd.Context(), items)
	if err != nil {
		return err
	}

	if selected.Header == "" {
		return nil
	}

	chosen, ok := entryOf[selected.Header]
	if !ok {
		return notifyLaunchError(n, "selected app not found")
	}

	bump(state, chosen.ID)
	_ = fsys.SetJSON(fsys.DataPath("vlx", "launcher"), "state", state)

	cmdStr := stripExecCodes(chosen.Exec)
	execCmd := exec.CommandContext(cmd.Context(), "sh", "-c", cmdStr)
	execCmd.Stdout = logs.Stdout()
	execCmd.Stderr = logs.Stderr()

	if err := execCmd.Start(); err != nil {
		return notifyLaunchError(n, fmt.Sprintf("failed to launch %s", chosen.Name))
	}

	return nil
}

// notifyLaunchError sends a notification and returns an error.
func notifyLaunchError(n *notify.Notify, msg string) error {
	_ = n.Send(msg, &notify.Details{
		Title:   "VLX Launcher",
		Urgency: "normal",
	})

	return fmt.Errorf("launcher: %s", msg)
}

// stripExecCodes removes desktop entry field codes from the Exec line.
func stripExecCodes(s string) string {
	codes := []string{"%f", "%F", "%u", "%U", "%d", "%D", "%n", "%N", "%i", "%c", "%k", "%v", "%m"}
	for _, code := range codes {
		s = strings.ReplaceAll(s, code, "")
	}

	return strings.TrimSpace(s)
}

// bump increments the usage count for an application.
func bump(s *State, id string) {
	u := s.Usage[id]
	u.Count++
	u.LastUsed = time.Now().Unix()
	s.Usage[id] = u
}

// rank returns entries sorted by the state rules.
func rank(entries []Entry, state *State) []Entry {
	ranked := make([]Entry, len(entries))
	copy(ranked, entries)

	sort.SliceStable(ranked, func(i, j int) bool {
		return ranked[i].Name < ranked[j].Name
	})

	sort.SliceStable(ranked, func(i, j int) bool {
		return state.Usage[ranked[i].ID].Count > state.Usage[ranked[j].ID].Count
	})

	return ranked
}
