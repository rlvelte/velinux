package packages

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/rlvelte/velinux/vlx/internal/core/guard"
	"github.com/rlvelte/velinux/vlx/internal/visuals/picker"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
	"github.com/spf13/cobra"
)

// TODO: REWORK FOR - quickshell, cache?, notify, progress, codestyle

// setup validates all requirements for further processing.
func setup(cmd *cobra.Command, _ []string) error {
	if err := errors.Join(guard.Network(), guard.Binaries("zypper", "fzf")); err != nil {
		return err
	}

	cmd.SetContext(context.WithValue(cmd.Context(), printer.ContextKey, printer.New()))
	cmd.SetContext(context.WithValue(cmd.Context(), picker.ContextKey, picker.New()))
	return nil
}

// Command returns the cobra command tree for vlx packages.
func Command() *cobra.Command {
	root := &cobra.Command{
		Use:               "packages",
		Short:             "Horribly bad packages installer",
		Long:              "Package install wrapper around zypper with interactive search.",
		PersistentPreRunE: setup,
		Args:              cobra.ArbitraryArgs,
		Aliases:           []string{"pgk", "pkg"}, // typo protection
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdInstall(cmd, strings.Join(args, " "))
		},
	}

	return root
}

// cmdInstall runs a search and then an info + install on the selection
func cmdInstall(cmd *cobra.Command, query string) error {
	pi := cmd.Context().Value(picker.ContextKey).(*picker.Picker)
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)

	client := &Zypper{}

	pkgs, err := client.Search(cmd.Context(), query)
	if err != nil {
		p.Error("Search failed: " + err.Error())
		return err
	}

	if len(pkgs) == 0 {
		p.Info("No packages found")
		return nil
	}

	latestPkgs := latest(pkgs)
	items := make([]picker.Item, len(latestPkgs))
	for i, pkg := range latestPkgs {
		desc := pkg.Description
		if desc == "" {
			desc = string(pkg.Type)
		}
		items[i] = picker.Item{
			Header:      pkg.Name,
			Description: desc,
		}
	}

	selected, err := pi.Select(cmd.Context(), items)
	if err != nil {
		return fmt.Errorf("fzf selection failed: %w", err)
	}

	pkgName := selected.Header
	info, err := client.Info(cmd.Context(), pkgName)
	if err != nil {
		p.Warn(fmt.Sprintf("Failed to get info for %s: %v", pkgName, err))
	} else {
		p.Info(info)
	}

	if !p.Confirm("Install selected packages?", true) {
		return nil
	}

	return client.Install(cmd.Context(), []string{pkgName})
}

// latest returns the current version of a packages.
func latest(pkgs []Package) []Package {
	latest := make(map[string]Package)
	for _, pkg := range pkgs {
		existing, exists := latest[pkg.Name]
		if !exists || pkg.Version > existing.Version {
			latest[pkg.Name] = pkg
		}
	}

	result := make([]Package, 0, len(latest))
	for _, pkg := range latest {
		result = append(result, pkg)
	}

	return result
}
