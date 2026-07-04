package hardware

import (
	"context"
	"fmt"

	"github.com/rlvelte/velinux/vlx/internal/core/printer"
	"github.com/spf13/cobra"
)

// setup validates all requirements for further processing.
func setup(cmd *cobra.Command, _ []string) error {
	cmd.SetContext(context.WithValue(cmd.Context(), printer.ContextKey, printer.New()))
	return nil
}

// Command returns the cobra command tree for vlx commandcenter hardware.
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "hardware [source]",
		Short:             "Horribly bad hardware monitor",
		Long:              "Query hardware information by source. Run without arguments to list available sources.",
		Aliases:           []string{"hw"},
		Args:              cobra.MaximumNArgs(1),
		PersistentPreRunE: setup,
		RunE: func(cmd *cobra.Command, args []string) error {
			p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)

			if len(args) == 0 {
				ListSources(p)
				return nil
			}

			jsonFlag, _ := cmd.Flags().GetBool("json")
			s := Find(args[0])
			if s == nil {
				return fmt.Errorf("unknown hardware source: %s", args[0])
			}

			return s.Run(cmd.Context(), p, jsonFlag)
		},
	}

	cmd.Flags().BoolP("json", "j", false, "output as JSON")
	return cmd
}
