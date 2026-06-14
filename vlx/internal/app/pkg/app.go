package pkg

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/rlvelte/velinux/vlx/internal/core/guard"
	"github.com/rlvelte/velinux/vlx/internal/core/picker"
	"github.com/rlvelte/velinux/vlx/internal/core/printer"
	"github.com/spf13/cobra"
)

// setup validates all requirements for further processing.
func setup(cmd *cobra.Command, _ []string) error {
	if err := errors.Join(guard.Connection(), guard.Binaries("zypper", "fzf")); err != nil {
		return err
	}

	cmd.SetContext(context.WithValue(cmd.Context(), printer.ContextKey, printer.New()))
	return nil
}

// Command returns the cobra command tree for vlx pkg.
func Command() *cobra.Command {
	root := &cobra.Command{
		Use:               "pkg",
		Short:             "Horribly bad package installer",
		Long:              "Package install wrapper around zypper with interactive search.",
		PersistentPreRunE: setup,
		Args:              cobra.ArbitraryArgs,
		Aliases:           []string{"pgk"}, // typo protection
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdInstall(cmd, strings.Join(args, " "))
		},
	}

	return root
}

// cmdInstall runs a search and then an info + install on the selection
func cmdInstall(cmd *cobra.Command, query string) error {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)

	p.Info("Searching packages...")
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
	items := make([]string, len(latestPkgs))
	for i, pkg := range latestPkgs {
		items[i] = format(pkg)
	}

	sort.Slice(items, func(i, j int) bool {
		return len(items[i]) < len(items[j])
	})

	fzf := picker.New().ForceFzf()
	selected, err := fzf.Select(cmd.Context(), items)
	if err != nil {
		return fmt.Errorf("fzf selection failed: %w", err)
	}

	if selected == "" {
		return nil
	}

	pkgName := extractName(selected)
	info, err := client.Info(cmd.Context(), pkgName)
	if err != nil {
		p.Warn(fmt.Sprintf("Failed to get info for %s: %v", pkgName, err))
	} else {
		p.Info(info)
	}

	if !p.Confirm("Install selected package?", true) {
		return nil
	}

	return client.Install(cmd.Context(), []string{pkgName})
}

// latest returns the current version of a package.
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

// format returns the fzf line style.
func format(pkg Package) string {
	status := " "
	if pkg.Installed {
		status = "i"
	}

	if pkg.Upgradable {
		status = "i+"
	}

	return fmt.Sprintf("%s | %s", status, pkg.Name)
}

// extractName returns the name of the package.
func extractName(item string) string {
	parts := strings.SplitN(item, "|", 2)
	if len(parts) > 1 {
		return strings.TrimSpace(parts[1])
	}

	return strings.TrimSpace(item)
}
