package hardware

import (
	"fmt"

	"github.com/rlvelte/velinux/vlx/internal/app/commandcenter/hardware/sources"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
	"github.com/spf13/cobra"
)

// Command returns the cobra command tree for vlx commandcenter hardware.
func Command() *cobra.Command {
	root := &cobra.Command{
		Use:     "hardware [source]",
		Short:   "Horribly bad hardware monitor",
		Long:    "Query hardware information by source. Run without arguments to list available sources.",
		Aliases: []string{"hw"},
		Args:    cobra.MaximumNArgs(1),
		RunE:    cmdRun,
	}

	listCmd := &cobra.Command{
		Use:     "list",
		Short:   "List available feed sources",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE:    cmdList,
	}

	root.AddCommand(listCmd)
	return root
}

func cmdList(cmd *cobra.Command, args []string) error {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)
	sources.List(p)

	return nil
}

func cmdRun(cmd *cobra.Command, args []string) error {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)

	s := sources.Find(args[0])
	if s == nil {
		return fmt.Errorf("unknown hardware source: %s", args[0])
	}

	return s.Run(p)
}
