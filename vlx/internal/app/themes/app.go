package themes

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mattn/go-isatty"
	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	"github.com/rlvelte/velinux/vlx/internal/visuals/notify"
	"github.com/rlvelte/velinux/vlx/internal/visuals/picker"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
	"github.com/spf13/cobra"
)

// TODO: REWORK FOR - progress

func Command() *cobra.Command {
	root := &cobra.Command{
		Use:          "themes",
		Short:        "Horribly bad theming manager",
		Long:         "Manage and switch between theme profiles.",
		Aliases:      []string{"theme"},
		Args:         cobra.NoArgs,
		SilenceUsage: !isatty.IsTerminal(os.Stdout.Fd()),
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
			Run:     cmdList,
		},
		&cobra.Command{
			Use:          "apply [theme]",
			Short:        "Apply a theme",
			Long:         "Apply a theme by name or interactively select from a list.",
			Aliases:      []string{"a", "ap"},
			Args:         cobra.MaximumNArgs(1),
			SilenceUsage: true,
			Run:          cmdApply,
		},
	)

	return root
}

// cmdList lists all available themes.
func cmdList(cmd *cobra.Command, _ []string) {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)

	dir := fsys.ConfigPath("vlx", "themes")
	active := current()
	themes, err := fsys.ListJSON(dir, decodeTheme)
	if err != nil {
		p.Error(err)
		return
	}

	if len(themes) == 0 {
		p.Success("No themes found")
		return
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
}

// cmdApply applies a selected theme.
func cmdApply(cmd *cobra.Command, args []string) {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)
	n := cmd.Context().Value(notify.ContextKey).(*notify.Notify)

	dir := fsys.ConfigPath("vlx", "themes")

	var fileName string
	if len(args) == 0 {
		fn, err := pick(cmd, dir)
		if err != nil {
			p.Error(err)
			return
		}

		fileName = fn
	} else {
		fileName = args[0]
	}

	theme, err := fsys.GetJSON(dir, fileName, decodeTheme)
	if err != nil {
		p.Error(err)
		return
	}

	currentPath := filepath.Join(dir, "current.json")
	if err := os.Remove(currentPath); err != nil && !os.IsNotExist(err) {
		p.Error(err)
		return
	}

	if err := os.Symlink(theme.filename+".json", currentPath); err != nil {
		p.Error(err)
		return
	}

	wallpaperPath := filepath.Join(dir, "wallpaper.png")
	if err := os.Remove(wallpaperPath); err != nil && !os.IsNotExist(err) {
		p.Error(err)
		return
	}

	if err := os.Symlink(theme.Wallpaper, wallpaperPath); err != nil {
		p.Error(err)
		return
	}

	if err := GenerateAll(*theme); err != nil {
		p.Error(err)
		return
	}

	_ = exec.Command("swaymsg", "reload").Run()
	_ = exec.Command("hyprctl", "reload").Run()
	_ = exec.Command("makoctl", "reload").Run()
	_ = exec.Command("mmsg", "dispatch", "reload_config").Run()

	_ = n.Send("Switched to theme "+theme.Name, &notify.Details{
		Title:   "VLX Themes",
		Urgency: "normal",
	})

	p.Success("Applied theme " + theme.Name)
	return
}

// pick selects a theme via an interactive picker.
func pick(cmd *cobra.Command, dir string) (string, error) {
	pkr, err := picker.New()
	if err != nil {
		return "", err
	}

	themes, err := fsys.ListJSON(dir, decodeTheme)
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
			Icon:        filepath.Join(dir, t.Logo),
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
	t, err := fsys.GetJSON(fsys.ConfigPath("vlx", "themes"), "current", decodeTheme)
	if err != nil {
		return ""
	}

	return t.filename
}
