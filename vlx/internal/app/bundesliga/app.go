package bundesliga

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	"github.com/rlvelte/velinux/vlx/internal/core/guard"
	vlxhttp "github.com/rlvelte/velinux/vlx/internal/core/http"
	"github.com/rlvelte/velinux/vlx/internal/core/printer"
	"github.com/spf13/cobra"
)

const apiBase = "https://api.openligadb.de"

func season() int {
	now := time.Now()
	if now.Month() >= time.August {
		return now.Year()
	}
	return now.Year() - 1
}

// setup validates preconditions and injects shared services into context.
func setup(cmd *cobra.Command, _ []string) error {
	if err := errors.Join(guard.Connection()); err != nil {
		return err
	}

	cmd.SetContext(context.WithValue(cmd.Context(), printer.ContextKey, printer.New()))
	cmd.SetContext(context.WithValue(cmd.Context(), vlxhttp.ContextKey, &vlxhttp.HTTP{}))
	return nil
}

// Command returns the cobra command tree for vlx bundesliga.
func Command() *cobra.Command {
	root := &cobra.Command{
		Use:               "bundesliga",
		Short:             "Horribly bad bundesliga tracker",
		Aliases:           []string{"bl"},
		PersistentPreRunE: setup,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	watchCmd := &cobra.Command{
		Use:   "setup",
		Short: "Setup notifications for ongoing matches",
		Long:  "Notify about score changes when selected team is playing",
		Args:  cobra.NoArgs,
		RunE:  cmdSetup,
	}

	watchCmd.Flags().Bool("poll", false, "")
	_ = watchCmd.Flags().MarkHidden("poll")

	root.AddCommand(watchCmd)
	root.AddCommand(&cobra.Command{
		Use:   "disable",
		Short: "Disable match tracking",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)
			if err := exec.Command("systemctl", "--user", "disable", "--now", "vlx-bl-tracker.timer").Run(); err != nil {
				return fmt.Errorf("disabling timer: %w", err)
			}

			p.Info("Tracker disabled")
			return nil
		},
	})
	return root
}

// cmdSetup lets the user set a favorite team and enables notifications
func cmdSetup(cmd *cobra.Command, _ []string) error {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)

	isPoll, _ := cmd.Flags().GetBool("poll")
	if isPoll {
		return poll(cmd.Context().Value(vlxhttp.ContextKey).(*vlxhttp.HTTP))
	}

	http := cmd.Context().Value(vlxhttp.ContextKey).(*vlxhttp.HTTP)

	store := fsys.NewStore(fsys.ConfigPath("vlx", "bundesliga"), decodeConfig, ".json")
	cfg, _ := store.Get("config")

	if err := selectTeam(cmd.Context(), http, &cfg); err != nil {
		return err
	}

	if err := saveConfig(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	if err := exec.Command("systemctl", "--user", "enable", "--now", "vlx-bl-tracker.timer").Run(); err != nil {
		return fmt.Errorf("enabling timer: %w", err)
	}

	p.Info(fmt.Sprintf("Tracker enabled for %s", cfg.Team.TeamName))
	return nil
}
