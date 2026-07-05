package fetch

import (
	"context"

	"github.com/rlvelte/velinux/vlx/internal/app/fetch/bundesliga"
	"github.com/rlvelte/velinux/vlx/internal/app/fetch/feeds"
	"github.com/rlvelte/velinux/vlx/internal/visuals/notify"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
	"github.com/spf13/cobra"
)

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

// Command returns the cobra command tree for vlx info.
func Command() *cobra.Command {
	root := &cobra.Command{
		Use:               "fetch",
		Short:             "Horribly bad information fetching",
		Long:              "Fetches information from web sources.",
		Aliases:           []string{"f", "info", "cc"},
		PersistentPreRunE: setup,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	root.PersistentFlags().BoolP("json", "j", false, "output as JSON")

	root.AddCommand(
		bundesliga.Command(),
		feeds.Command(),
	)

	return root
}
