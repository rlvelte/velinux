package fetch

import (
	"github.com/rlvelte/velinux/vlx/internal/app/fetch/bundesliga"
	"github.com/rlvelte/velinux/vlx/internal/app/fetch/feeds"
	"github.com/rlvelte/velinux/vlx/internal/app/fetch/hardware"
	"github.com/spf13/cobra"
)

// Command returns the cobra command tree for vlx info.
func Command() *cobra.Command {
	root := &cobra.Command{
		Use:     "fetch",
		Short:   "Horribly bad information fetching",
		Long:    "Fetches information from web and local sources.",
		Aliases: []string{"f"},
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
