package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	"github.com/rlvelte/velinux/vlx/internal/visuals/notify"
	"github.com/rlvelte/velinux/vlx/internal/visuals/picker"
	"github.com/spf13/cobra"
)

// TODO: REWORK FOR - quickshell performance, cache?, codestyle

func Command() *cobra.Command {
	return &cobra.Command{
		Use:     "launcher",
		Short:   "Horribly bad application launcher",
		Long:    "Scan desktop entries and launch via the quickshell/fzf picker.",
		Aliases: []string{"run", "app"},
		Args:    cobra.NoArgs,
		RunE:    cmdLaunch,
	}
}

func cmdLaunch(cmd *cobra.Command, _ []string) error {
	n := cmd.Context().Value(notify.ContextKey).(*notify.Notify)

	store := fsys.NewStore(fsys.ConfigPath("vlx", "launcher"), decodeConfig, ".json")
	cfg, err := store.Get("config")
	if err != nil {
		return fmt.Errorf("launcher config: %w", err)
	}

	stateStore := fsys.NewStore(fsys.DataPath("vlx", "launcher"), decodeState, ".json")
	state, err := stateStore.Get("state")
	if err != nil {
		return fmt.Errorf("launcher state: %w", err)
	}

	entries := Scan()

	var visible []Entry
	for _, e := range entries {
		if e.NoDisplay || e.Hidden || IsIgnored(cfg, e.ID) {
			continue
		}

		visible = append(visible, e)
	}

	ranked := Rank(visible, state)

	items := make([]picker.Item, len(ranked))
	for i, e := range ranked {
		items[i] = picker.Item{
			Icon:        resolveIcon(e.Icon),
			Header:      e.Name,
			Description: e.Comment,
		}
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

	var chosen *Entry
	for _, e := range visible {
		if e.Name == selected.Header {
			chosen = &e
			break
		}
	}

	if chosen == nil {
		fmtErr := fmt.Errorf("selected app not found")
		_ = n.Send(fmtErr.Error(), &notify.Details{
			Title:   "Couldn't launch",
			Urgency: "normal",
		})

		return fmtErr
	}

	Bump(state, chosen.ID)
	_ = SaveState(state)

	cmdStr := stripExecCodes(chosen.Exec)
	execCmd := exec.CommandContext(cmd.Context(), "sh", "-c", cmdStr)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	if err := execCmd.Start(); err != nil {
		fmtErr := fmt.Errorf("failed to launch %s: %w", chosen.Name, err)
		_ = n.Send(fmtErr.Error(), &notify.Details{
			Title:   "Couldn't launch",
			Urgency: "normal",
		})

		return fmtErr
	}

	return nil
}

// stripExecCodes removes desktop entry field codes from the Exec line.
func stripExecCodes(s string) string {
	codes := []string{"%f", "%F", "%u", "%U", "%d", "%D", "%n", "%N", "%i", "%c", "%k", "%v", "%m"}
	for _, code := range codes {
		s = strings.ReplaceAll(s, code, "")
	}
	return strings.TrimSpace(s)
}
