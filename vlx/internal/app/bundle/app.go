package bundle

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	"github.com/rlvelte/velinux/vlx/internal/visuals/picker"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
	"github.com/spf13/cobra"
)

// TODO: REWORK FOR - notify, progress

// Command returns the cobra command tree for vlx bundle.
func Command() *cobra.Command {
	root := &cobra.Command{
		Use:     "bundle",
		Short:   "Horribly bad bundle installer",
		Long:    "Install/Compile predefined recipes with shell hooks.",
		Aliases: []string{"bun"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	root.AddCommand(
		&cobra.Command{
			Use:     "list",
			Short:   "List available bundles",
			Long:    "List all available bundles with their segments.",
			Aliases: []string{"l", "ls"},
			Args:    cobra.NoArgs,
			RunE:    cmdList,
		},

		&cobra.Command{
			Use:     "install [bundle]",
			Short:   "Install a bundle",
			Long:    "Install a bundle by name or interactively select from a list.",
			Aliases: []string{"i", "in"},
			Args:    cobra.MaximumNArgs(1),
			RunE:    cmdInstall,
		},
	)

	return root
}

// cmdList lists out all bundles and peek at its contents.
func cmdList(cmd *cobra.Command, _ []string) error {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)

	bundlesDir := fsys.ConfigPath("vlx", "bundles")
	store := fsys.NewStore(bundlesDir, decodeBundle, ".json")

	bundles, err := store.List()
	if err != nil {
		p.Error(err.Error())
		return err
	}

	if len(bundles) == 0 {
		p.Print("No bundles found")
		return nil
	}

	headers := []string{"Name", "Zypper", "Flatpak", "Hooks"}
	rows := make([][]string, 0, len(bundles))
	for _, b := range bundles {
		hooks := "no"
		if len(b.PreHook) > 0 || len(b.PostHook) > 0 {
			hooks = "yes"
		}

		rows = append(rows, []string{
			b.Info.Name,
			strconv.Itoa(len(b.Zypper)),
			strconv.Itoa(len(b.Flatpak)),
			hooks,
		})
	}

	p.Table(headers, rows)
	return nil
}

// cmdInstall installs a selected bundle.
func cmdInstall(cmd *cobra.Command, args []string) error {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)

	bundlesDir := fsys.ConfigPath("vlx", "bundles")
	store := fsys.NewStore(bundlesDir, decodeBundle, ".json")

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

	bundle, err := store.Get(fileName)
	if err != nil {
		p.Error(err.Error())
		return err
	}

	steps := []struct {
		enabled bool
		run     func() error
	}{
		{len(bundle.Repos) > 0, func() error { return repo(bundle.Repos) }},
		{len(bundle.PreHook) > 0, func() error { return sh(strings.Join(bundle.PreHook, " && ")) }},
		{len(bundle.Zypper) > 0, func() error { return zypper(bundle.Zypper) }},
		{len(bundle.Flatpak) > 0, func() error { return flatpak(bundle.Flatpak) }},
		{len(bundle.PostHook) > 0, func() error { return sh(strings.Join(bundle.PostHook, " && ")) }},
	}

	for _, step := range steps {
		if !step.enabled {
			continue
		}

		if err := step.run(); err != nil {
			p.Error(err.Error())
			return err
		}
	}

	return nil
}

// pick selects a bundle via an interactive picker.
func pick(cmd *cobra.Command, store *fsys.Store[Bundle]) (string, error) {
	pkr, err := picker.New()
	if err != nil {
		return "", err
	}

	bundles, err := store.List()
	if err != nil {
		return "", err
	}

	items := make([]picker.Item, len(bundles))
	lookup := make(map[string]string, len(bundles))
	for i, b := range bundles {
		items[i] = picker.Item{
			Icon:        b.Info.Icon,
			Header:      b.Info.Name,
			Description: b.Info.Description,
		}

		lookup[b.Info.Name] = b.FileName
	}

	selected, err := pkr.Select(cmd.Context(), items)
	if err != nil {
		return "", err
	}

	return lookup[selected.Header], nil
}

// sh executes a cmd with basic shell
func sh(cmdStr string) error {
	cmd := exec.Command("sh", "-c", cmdStr)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// zypper runs a simple install
func zypper(pkgs []string) error {
	args := append([]string{"zypper", "install", "-y"}, pkgs...)
	cmd := exec.Command("sudo", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// flatpak runs a simple install
func flatpak(pkgs []string) error {
	args := append([]string{"flatpak", "install", "-y"}, pkgs...)
	cmd := exec.Command("sudo", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// repo adds a repository to zypper
func repo(repos []Repo) error {
	for _, repo := range repos {
		cmd := exec.Command("sudo", "zypper", "ar", repo.URL, repo.Alias)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add repo %q: %w", repo.Alias, err)
		}
	}

	return nil
}
