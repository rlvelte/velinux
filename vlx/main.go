package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/rvelte/vlx/internal/app/pkg"
	"github.com/rvelte/vlx/internal/app/stat"
	"github.com/rvelte/vlx/internal/app/themes"
)

func main() {
	root := &cobra.Command{
		Use:   "vlx",
		Short: "vlx utility application by rvelte",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	root.AddCommand(
		completions(),
		pkg.Command(),
		stat.Command(),
		themes.Command(),
	)

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}

func completions() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish]",
		Short: "Generate shell completion script",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			}

			return fmt.Errorf("unsupported shell: %s", args[0])
		},
	}
}
