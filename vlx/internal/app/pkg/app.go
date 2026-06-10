package pkg

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rvelte/vlx/internal/core/picker"
	"github.com/rvelte/vlx/internal/core/printer"
	"github.com/rvelte/vlx/internal/system/guard"
	"github.com/rvelte/vlx/internal/system/xdg"
	"github.com/spf13/cobra"
)

// stateDir returns the location of system-wide installed pkg schemes.
func stateDir() string {
	return xdg.ConfigPath("vlx", "pkg")
}

// ensure validates all requirements for further processing.
func ensure(_ *cobra.Command, _ []string) error {
	return errors.Join(guard.OS(), guard.Connection())
}

// Command returns the cobra command tree for vlx pkg.
func Command() *cobra.Command {
	root := &cobra.Command{
		Use:               "pkg",
		Short:             "Opinionated package scheme installer",
		PersistentPreRunE: ensure,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	root.AddCommand(
		&cobra.Command{
			Use:     "install [scheme]",
			Short:   "Install a scheme",
			Aliases: []string{"in"},
			Args:    cobra.MaximumNArgs(1),
			RunE:    cmdInstall,
		},
		&cobra.Command{
			Use:     "list",
			Short:   "List all available schemes",
			Aliases: []string{"ls"},
			RunE:    cmdList,
		},
	)

	return root
}

func cmdInstall(cmd *cobra.Command, args []string) error {
	var name string

	if len(args) == 0 {
		schemes, err := List()
		if err != nil {
			return fmt.Errorf("failed to list schemes: %w", err)
		}

		if len(schemes) == 0 {
			return errors.New("no schemes available")
		}

		items := make([]string, len(schemes))
		for i, s := range schemes {
			items[i] = s.Name
		}

		p := picker.New()
		if p == nil {
			return errors.New("no picker backend available")
		}

		selected, err := p.Select(cmd.Context(), items)
		if err != nil {
			return fmt.Errorf("scheme selection cancelled: %w", err)
		}

		name = selected
	} else {
		name = args[0]
	}

	s, err := Load(name)
	if err != nil {
		return err
	}

	return Install(cmd.Context(), s)
}

func cmdList(_ *cobra.Command, _ []string) error {
	p := printer.New()

	schemes, err := List()
	if err != nil {
		return err
	}

	if len(schemes) == 0 {
		p.Info(fmt.Sprintf("No schemes found in %s", stateDir()))
		return nil
	}

	var rows [][]string
	for _, s := range schemes {
		total := len(s.Repos) + len(s.Zypper) + len(s.Flatpak)

		hooks := ""
		if len(s.PreHook) > 0 || len(s.PostHook) > 0 {
			var h []string
			if len(s.PreHook) > 0 {
				h = append(h, "pre")
			}
			if len(s.PostHook) > 0 {
				h = append(h, "post")
			}
			hooks = strings.Join(h, ", ")
		}

		rows = append(rows, []string{s.Name, fmt.Sprintf("%d", total), hooks})
	}

	p.Table([]string{"Scheme", "Packages", "Hooks"}, rows)
	return nil
}
