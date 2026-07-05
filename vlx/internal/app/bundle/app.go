package bundle

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	"github.com/rlvelte/velinux/vlx/internal/core/guard"
	"github.com/rlvelte/velinux/vlx/internal/visuals/picker"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
	"github.com/spf13/cobra"
)

// TODO: REWORK FOR - quickshell, notify, progress, codestyle

// setup validates all requirements for further processing.
func setup(cmd *cobra.Command, _ []string) error {
	if err := errors.Join(guard.Network(), guard.Binaries("bash", "fzf")); err != nil {
		return err
	}

	cmd.SetContext(context.WithValue(cmd.Context(), printer.ContextKey, printer.New()))
	return nil
}

// Command returns the cobra command tree for vlx bundle.
func Command() *cobra.Command {
	root := &cobra.Command{
		Use:               "bundle",
		Short:             "Horribly bad bundle installer",
		Long:              "Install predefined recipes with shell hooks.",
		PersistentPreRunE: setup,
		Aliases:           []string{"bun"},
		Args:              cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	root.AddCommand(
		&cobra.Command{
			Use:     "list",
			Short:   "List available bundles",
			Long:    "List all available bundles with their segments.",
			Aliases: []string{"ls"},
			Args:    cobra.NoArgs,
			RunE:    cmdList,
		},
		&cobra.Command{
			Use:     "install [bundle]",
			Short:   "Install a bundle",
			Long:    "Install a bundle by name or interactively select from a list.",
			Aliases: []string{"in"},
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
		return err
	}

	if len(bundles) == 0 {
		p.Info("No bundles found")
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
			b.Name,
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

	var bundleName string

	if len(args) == 0 {
		bundles, err := store.List()
		if err != nil {
			return err
		}

		sort.Slice(bundles, func(i, j int) bool {
			return bundles[i].Name < bundles[j].Name
		})

		items := make([]picker.Item, len(bundles))
		for i, b := range bundles {
			desc := fmt.Sprintf("%d zypper, %d flatpak", len(b.Zypper), len(b.Flatpak))
			items[i] = picker.Item{
				Header:      b.Name,
				Description: desc,
			}
		}

		pkr := picker.New().ForceFzf()
		if pkr == nil {
			return fmt.Errorf("no picker available")
		}

		selected, err := pkr.Select(cmd.Context(), items)
		if err != nil {
			return err
		}

		bundleName = selected.Header
	} else {
		bundleName = args[0]
	}

	bundle, err := store.Get(bundleName)
	if err != nil {
		return fmt.Errorf("bundle %q not found", bundleName)
	}

	if len(bundle.Repos) > 0 {
		if err := p.Spinner("\nAdding repositories", func() error {
			return repo(bundle.Repos)
		}); err != nil {
			return fmt.Errorf("adding repos failed: %w", err)
		}
	}

	if len(bundle.PreHook) > 0 {
		combined := strings.Join(bundle.PreHook, " && ")
		if err := p.Spinner("\nRunning pre-hooks", func() error {
			return sh(combined)
		}); err != nil {
			return fmt.Errorf("pre-hooks failed: %w", err)
		}
	}

	if len(bundle.Zypper) > 0 {
		if err := p.Spinner("\nInstalling zypper packages", func() error {
			return zypper(bundle.Zypper)
		}); err != nil {
			return fmt.Errorf("zypper install failed: %w", err)
		}
	}

	if len(bundle.Flatpak) > 0 {
		if err := p.Spinner("\nInstalling flatpak packages", func() error {
			return flatpak(bundle.Flatpak)
		}); err != nil {
			return fmt.Errorf("flatpak install failed: %w", err)
		}
	}

	if len(bundle.PostHook) > 0 {
		combined := strings.Join(bundle.PostHook, " && ")
		if err := p.Spinner("\nRunning post-hooks", func() error {
			return sh(combined)
		}); err != nil {
			return fmt.Errorf("post-hooks failed: %w", err)
		}
	}

	p.Info(fmt.Sprintf("Bundle %q installed", bundleName))
	return nil
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
