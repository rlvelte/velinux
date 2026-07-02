package themes

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	"github.com/rlvelte/velinux/vlx/internal/core/notify"
	"github.com/rlvelte/velinux/vlx/internal/core/picker"
	"github.com/rlvelte/velinux/vlx/internal/core/printer"
	"github.com/spf13/cobra"
)

// setup configures all requirements and guards against wrong usage.
func setup(cmd *cobra.Command, _ []string) error {
	cmd.SetContext(context.WithValue(cmd.Context(), notify.ContextKey, notify.New()))
	cmd.SetContext(context.WithValue(cmd.Context(), printer.ContextKey, printer.New()))
	return nil
}

func Command() *cobra.Command {
	root := &cobra.Command{
		Use:               "themes",
		Short:             "Horribly bad theming manager for velinux",
		Long:              "Manage and switch between theme profiles for velinux.",
		PersistentPreRunE: setup,
		Args:              cobra.NoArgs,
		Aliases:           []string{"theme"},
		SilenceUsage:      true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmdListCmd := &cobra.Command{
		Use:     "list",
		Short:   "List available theme profiles",
		Long:    "List all available theme profiles with their icons and IDs.",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE:    cmdList,
	}
	cmdListCmd.Flags().BoolP("json", "j", false, "output as JSON")

	root.AddCommand(
		cmdListCmd,
		&cobra.Command{
			Use:     "apply [theme]",
			Short:   "Apply a theme",
			Long:    "Apply a theme by name or interactively select from a list.",
			Aliases: []string{"sw"},
			Args:    cobra.MaximumNArgs(1),
			RunE:    cmdApply,
		},
	)

	return root
}

type themeJSON struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Active bool   `json:"active"`
}

func cmdList(cmd *cobra.Command, _ []string) error {
	themesDir := fsys.ConfigPath("vlx", "themes")
	store := fsys.NewStore(themesDir, decodeTheme, ".conf")
	active := current()

	all, err := store.List()
	if err != nil {
		return err
	}

	seen := make(map[string]bool)
	var list []*Theme
	for _, t := range all {
		if seen[t.Id] {
			continue
		}
		if filepath.Base(t.Path) == "current.conf" {
			continue
		}

		seen[t.Id] = true
		list = append(list, t)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	jsonFlag, _ := cmd.Flags().GetBool("json")
	if jsonFlag {
		var out []themeJSON
		for _, t := range list {
			out = append(out, themeJSON{
				Id:     t.Id,
				Name:   t.Name,
				Icon:   t.Icon,
				Active: t.Id == active,
			})
		}

		data, err := json.Marshal(out)
		if err != nil {
			return err
		}

		fmt.Println(string(data))
		return nil
	}

	p, _ := cmd.Context().Value(printer.ContextKey).(*printer.Printer)
	headers := []string{"ACTIVE", "ID", "Name"}
	var rows [][]string
	for _, t := range list {
		marker := ""
		if t.Id == active {
			marker = "*"
		}

		rows = append(rows, []string{marker, t.Id, t.Name})
	}

	p.Table(headers, rows)
	return nil
}

func cmdApply(cmd *cobra.Command, args []string) error {
	themesDir := fsys.ConfigPath("vlx", "themes")

	store := fsys.NewStore(themesDir, decodeTheme, ".conf")
	all, err := store.List()
	if err != nil {
		return err
	}

	seen := make(map[string]bool)
	var themes []*Theme
	for _, t := range all {
		if seen[t.Id] {
			continue
		}
		if filepath.Base(t.Path) == "current.conf" {
			continue
		}
		seen[t.Id] = true
		themes = append(themes, t)
	}

	var theme *Theme
	if len(args) == 0 {
		pkr := picker.New()
		if pkr == nil {
			return fmt.Errorf("no picker available")
		}

		sort.Slice(themes, func(i, j int) bool {
			return themes[i].Name < themes[j].Name
		})

		names := make([]string, len(themes))
		for i, t := range themes {
			names[i] = t.Name
		}

		selected, err := pkr.Select(cmd.Context(), names)
		if err != nil {
			return err
		}

		for _, t := range themes {
			if t.Name == selected {
				theme = t
				break
			}
		}
	} else {
		req := args[0]
		for _, t := range themes {
			if t.Id == req || t.Name == req {
				theme = t
				break
			}
		}
	}

	if theme == nil {
		return fmt.Errorf("theme not found")
	}

	data, err := os.ReadFile(theme.Path)
	if err != nil {
		return err
	}

	content, err := decodeThemeContent("", theme.Path, data)
	if err != nil {
		return err
	}

	currentPath := filepath.Join(themesDir, "current.conf")
	if err := os.Remove(currentPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.Symlink(filepath.Base(theme.Path), currentPath); err != nil {
		return err
	}

	wallpaperPath := filepath.Join(themesDir, "current.png")
	if err := os.Remove(wallpaperPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.Symlink(theme.Wallpaper, wallpaperPath); err != nil {
		return err
	}

	if err := GenerateAll(*content); err != nil {
		return err
	}

	if err := exec.Command("swaymsg", "reload").Run(); err != nil {
		p, _ := cmd.Context().Value(printer.ContextKey).(*printer.Printer)
		if p != nil {
			p.Warn("sway reload failed (sway may not be running)")
		}
	}

	if err := exec.Command("hyprctl", "reload").Run(); err != nil {
		p, _ := cmd.Context().Value(printer.ContextKey).(*printer.Printer)
		if p != nil {
			p.Warn("hypr reload failed (Hyprland may not be running)")
		}
	}

	if err := exec.Command("makoctl", "reload").Run(); err != nil {
		p, _ := cmd.Context().Value(printer.ContextKey).(*printer.Printer)
		if p != nil {
			p.Warn("mako reload failed (mako may not be running)")
		}
	}

	if err := exec.Command("mmsg", "dispatch", "reload_config").Run(); err != nil {
		p, _ := cmd.Context().Value(printer.ContextKey).(*printer.Printer)
		if p != nil {
			p.Warn("mango reload failed (mango may not be running)")
		}
	}

	n := notify.New()
	_ = n.Send("Switched to theme "+theme.Name, &notify.Details{
		Title:   "VLX Themes",
		Urgency: "normal",
	})

	p, _ := cmd.Context().Value(printer.ContextKey).(*printer.Printer)
	if p != nil {
		p.Info("Applied theme " + theme.Name)
	}

	return nil
}

// current returns the currently active theme id.
func current() string {
	themes := fsys.ConfigPath("vlx", "themes")
	data, err := os.ReadFile(filepath.Join(themes, "current.conf"))
	if err != nil {
		return ""
	}

	t, err := decodeTheme("current", "", data)
	if err != nil {
		return ""
	}

	return t.Id
}
