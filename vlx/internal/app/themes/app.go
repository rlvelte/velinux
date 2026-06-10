package themes

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/rvelte/vlx/internal/core/notify"
	"github.com/rvelte/vlx/internal/core/picker"
	"github.com/rvelte/vlx/internal/system/guard"
	"github.com/rvelte/vlx/internal/system/xdg"
	"github.com/spf13/cobra"
)

// Theme information for identification.
type Theme struct {
	Icon     string // Icon in ASCII format
	Id       string // Id for this theme
	Name     string // Human-readable Name
	Location string // Filesystem Location
	Valid    bool   // If .ini file is Valid
}

// stateDir returns the location of system-wide installed themes.
func stateDir() string {
	return xdg.ConfigPath("vlx", "themes")
}

// ensure validates all requirements for further processing.
func ensure(_ *cobra.Command, _ []string) error {
	return guard.OS()
}

// Command returns the cobra command tree for vlx themes.
func Command() *cobra.Command {
	root := &cobra.Command{
		Use:               "themes",
		Short:             "Theming manager for velinux",
		PersistentPreRunE: ensure,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	root.AddCommand(
		&cobra.Command{
			Use:     "list",
			Short:   "List available theme profiles",
			Aliases: []string{"ls"},
			RunE:    cmdList,
		},
		&cobra.Command{
			Use:     "switch [theme]",
			Short:   "Switch to theme and apply",
			Aliases: []string{"sw"},
			RunE:    cmdSwitch,
		},
		&cobra.Command{
			Use:   "waybar",
			Short: "Show current theme in waybar format",
			RunE:  cmdWaybar,
		},
	)

	return root
}

// cmdList lists all available themes.
func cmdList(_ *cobra.Command, _ []string) error {
	themes, err := listThemes()
	if err != nil {
		return fmt.Errorf("failed to list themes: %w", err)
	}

	if len(themes) == 0 {
		slog.Warn("No themes found. Themes are located in: " + stateDir())
		return nil
	}

	for _, t := range themes {
		fmt.Printf("%s\t%s\t%s\n", t.Icon, t.Id, t.Name)
	}

	return nil
}

// cmdSwitch switches to the specified theme.
func cmdSwitch(cmd *cobra.Command, args []string) error {
	var name string

	if len(args) == 0 {
		themes, err := listThemes()
		if err != nil {
			return fmt.Errorf("failed to list themes: %w", err)
		}

		if len(themes) == 0 {
			return errors.New("no themes available")
		}

		items := make([]string, len(themes))
		for i, t := range themes {
			items[i] = t.Name
		}

		picker := picker.New()
		if picker == nil {
			return errors.New("no picker backend available")
		}

		selected, err := picker.Select(cmd.Context(), items)
		if err != nil {
			return fmt.Errorf("theme selection cancelled: %w", err)
		}

		name = selected
	} else {
		name = args[0]
	}

	if err := switchTheme(name); err != nil {
		return err
	}

	_ = notify.NotifyWith(&notify.NotifyConfig{
		Title:   "vlx",
		Urgency: "normal",
	}, "You've switched to theme "+name)

	return nil
}

// cmdWaybar shows the current theme name in waybar format.
func cmdWaybar(_ *cobra.Command, _ []string) error {
	theme, err := currentTheme()
	if err != nil {
		return fmt.Errorf("failed to read current theme: %w", err)
	}

	if theme.Id == "suse" {
		fmt.Printf("%s ", theme.Icon)
		return nil
	}

	fmt.Printf("%s", theme.Icon)
	return nil
}
