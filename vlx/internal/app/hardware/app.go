package hardware

import (
	"context"
	"fmt"

	"github.com/rlvelte/velinux/vlx/internal/app/hardware/sources"
	"github.com/rlvelte/velinux/vlx/internal/visuals/notify"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
	"github.com/spf13/cobra"
)

// TODO: REWORK FOR - quickshell, notify, codestyle

// setup validates all requirements for further processing.
func setup(cmd *cobra.Command, _ []string) error {
	p := printer.New()
	if jsonFlag, _ := cmd.Flags().GetBool("json"); jsonFlag {
		p = p.ForceJSON()
	}

	cmd.SetContext(context.WithValue(cmd.Context(), printer.ContextKey, p))
	cmd.SetContext(context.WithValue(cmd.Context(), notify.ContextKey, notify.New()))
	return nil
}

// Command returns the cobra command tree for vlx hardware.
func Command() *cobra.Command {
	root := &cobra.Command{
		Use:               "hardware [source]",
		Short:             "Horribly bad hardware monitor",
		Long:              "Query hardware information by source.",
		Aliases:           []string{"hw"},
		Args:              cobra.MaximumNArgs(1),
		PersistentPreRunE: setup,
		RunE:              cmdRun,
	}

	listCmd := &cobra.Command{
		Use:     "list",
		Short:   "List available feed sources",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE:    cmdList,
	}

	root.PersistentFlags().BoolP("json", "j", false, "output as JSON")
	root.AddCommand(listCmd)
	return root
}

// cmdList lists out all the available hardware sources.
func cmdList(cmd *cobra.Command, _ []string) error {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)
	sources.List(p)

	return nil
}

// cmdRun fetches available information from the source.
func cmdRun(cmd *cobra.Command, args []string) error {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)

	s := sources.Find(args[0])
	if s == nil {
		return fmt.Errorf("unknown hardware source: %s", args[0])
	}

	return s.Run(p)
}
