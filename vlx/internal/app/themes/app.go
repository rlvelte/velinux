package themes

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	"github.com/rlvelte/velinux/vlx/internal/visuals/notify"
	"github.com/rlvelte/velinux/vlx/internal/visuals/picker"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
	"github.com/spf13/cobra"
)

// TODO: REWORK FOR - progress

func Command() *cobra.Command {
	root := &cobra.Command{
		Use:     "themes",
		Short:   "Horribly bad theming manager",
		Long:    "Manage and switch between theme profiles.",
		Aliases: []string{"theme"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	root.AddCommand(
		&cobra.Command{
			Use:     "list",
			Short:   "List available theme profiles",
			Long:    "List all available theme profiles with their icons and IDs.",
			Aliases: []string{"l", "ls"},
			Args:    cobra.NoArgs,
			RunE:    cmdList,
		},
		&cobra.Command{
			Use:     "apply [theme]",
			Short:   "Apply a theme",
			Long:    "Apply a theme by name or interactively select from a list.",
			Aliases: []string{"a", "ap"},
			Args:    cobra.MaximumNArgs(1),
			RunE:    cmdApply,
		},
	)

	return root
}

// cmdList lists all available themes.
func cmdList(cmd *cobra.Command, _ []string) error {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)

	store := fsys.NewStore(fsys.ConfigPath("vlx", "themes"), decodeTheme, ".json")

	active := current()
	themes, err := store.List()
	if err != nil {
		p.Error(err.Error())
		return err
	}

	if len(themes) == 0 {
		p.Print("No themes found")
		return nil
	}

	headers := []string{"Name", "Description", "Active"}
	var rows [][]string
	for _, t := range themes {
		marker := ""
		if t.filename == active {
			marker = "*"
		}

		rows = append(rows, []string{
			t.Name,
			t.Description,
			marker,
		})
	}

	p.Table(headers, rows)
	return nil
}

// cmdApply applies a selected theme.
func cmdApply(cmd *cobra.Command, args []string) error {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)
	n := cmd.Context().Value(notify.ContextKey).(*notify.Notify)

	store := fsys.NewStore(fsys.ConfigPath("vlx", "themes"), decodeTheme, ".json")

	var fileName string
	if len(args) == 0 {
		fn, err := pick(cmd, store)
		if err != nil {
			p.Error(err.Error())
			return err
		}

		fileName = fn
	} else {
		fileName = args[0]
	}

	theme, err := store.Get(fileName)
	if err != nil {
		p.Error(err.Error())
		return err
	}

	currentPath := filepath.Join(store.Dir, "current.json")
	if err := os.Remove(currentPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := os.Symlink(theme.filename+".json", currentPath); err != nil {
		return err
	}

	wallpaperPath := filepath.Join(store.Dir, "wallpaper.png")
	if err := os.Remove(wallpaperPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := os.Symlink(theme.Wallpaper, wallpaperPath); err != nil {
		return err
	}

	if err := GenerateAll(*theme); err != nil {
		return err
	}

	if err := exec.Command("swaymsg", "reload").Run(); err != nil {
		p.Warn("sway reload failed (sway may not be running)")
	}

	if err := exec.Command("hyprctl", "reload").Run(); err != nil {
		p.Warn("hypr reload failed (Hyprland may not be running)")
	}

	if err := exec.Command("makoctl", "reload").Run(); err != nil {
		p.Warn("mako reload failed (mako may not be running)")
	}

	if err := exec.Command("mmsg", "dispatch", "reload_config").Run(); err != nil {
		p.Warn("mango reload failed (mango may not be running)")
	}

	_ = n.Send("Switched to theme "+theme.Name, &notify.Details{
		Title:   "VLX Themes",
		Urgency: "normal",
	})

	p.Print("Applied theme " + theme.Name)
	return nil
}

// pick selects a theme via an interactive picker.
func pick(cmd *cobra.Command, store *fsys.Store[*Theme]) (string, error) {
	pkr, err := picker.New()
	if err != nil {
		return "", err
	}

	themes, err := store.List()
	if err != nil {
		return "", err
	}

	items := make([]picker.Item, 0, len(themes))
	lookup := make(map[string]string, len(themes))
	for _, t := range themes {
		if t.filename == "current" {
			continue
		}
		items = append(items, picker.Item{
			Icon:        filepath.Join(store.Dir, t.Logo),
			Header:      t.Name,
			Description: t.Description,
		})

		lookup[t.Name] = t.filename
	}

	selected, err := pkr.Select(cmd.Context(), items)
	if err != nil {
		return "", err
	}

	return lookup[selected.Header], nil
}

// current returns the currently active theme id.
func current() string {
	data, err := os.ReadFile(filepath.Join(fsys.ConfigPath("vlx", "themes"), "current.json"))
	if err != nil {
		return ""
	}

	t, err := decodeTheme("current", "", data)
	if err != nil {
		return ""
	}

	return t.filename
}
