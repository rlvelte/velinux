package hardware

import (
	"github.com/rlvelte/velinux/vlx/internal/app/fetch/hardware/sources"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
	"github.com/spf13/cobra"
)

// TODO: REWORK FOR - quickshell, notify, codestyle

// Command returns the cobra command tree for vlx hardware.
func Command() *cobra.Command {
	root := &cobra.Command{
		Use:     "hardware [source]",
		Short:   "Horribly bad hardware monitor",
		Long:    "Query hardware information by source.",
		Aliases: []string{"hw"},
		Args:    cobra.MaximumNArgs(1),
		RunE:    cmdRun,
	}

	root.AddCommand(
		&cobra.Command{
			Use:     "list",
			Short:   "List available feed sources",
			Aliases: []string{"l", "ls"},
			Args:    cobra.NoArgs,
			RunE:    cmdList,
		})

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

	s, err := sources.Find(args[0])
	if err != nil {
		p.Error(err.Error())
		return err
	}

	return s.Run(p)
}
