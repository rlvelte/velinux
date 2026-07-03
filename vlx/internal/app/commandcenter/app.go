package commandcenter

import (
	"context"

	"github.com/rlvelte/velinux/vlx/internal/app/commandcenter/bundesliga"
	"github.com/rlvelte/velinux/vlx/internal/app/commandcenter/feeds"
	"github.com/rlvelte/velinux/vlx/internal/app/commandcenter/hardware"
	"github.com/rlvelte/velinux/vlx/internal/core/printer"
	"github.com/spf13/cobra"
)

// setup validates all requirements for further processing.
func setup(cmd *cobra.Command, _ []string) error {
	cmd.SetContext(context.WithValue(cmd.Context(), printer.ContextKey, printer.New()))
	return nil
}

// Command returns the cobra command tree for vlx commandcenter.
func Command() *cobra.Command {
	root := &cobra.Command{
		Use:               "commandcenter",
		Short:             "Horribly bad command center",
		Long:              "System command center with hardware monitoring and more.",
		Aliases:           []string{"cc"},
		PersistentPreRunE: setup,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	root.AddCommand(
		hardware.Command(),
		bundesliga.Command(),
		feeds.Command(),
	)

	return root
}
