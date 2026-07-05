package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rlvelte/velinux/vlx/internal/app/bundle"
	"github.com/rlvelte/velinux/vlx/internal/app/fetch"
	"github.com/rlvelte/velinux/vlx/internal/app/launcher"
	"github.com/rlvelte/velinux/vlx/internal/app/mise"
	"github.com/rlvelte/velinux/vlx/internal/app/packages"
	"github.com/rlvelte/velinux/vlx/internal/app/themes"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "vlx",
		Short: "Horribly bad utility application",
		Long:  "VeLinux centered command utility by rvelte.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	root.AddCommand(
		completions(),
		packages.Command(),
		themes.Command(),
		bundle.Command(),
		fetch.Command(),
		mise.Command(),
		launcher.Command(),
	)

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}

func completions() *cobra.Command {
	return &cobra.Command{
		Use:    "completion [bash|zsh|fish]",
		Short:  "Generate shell completion script",
		Hidden: true,
		Args:   cobra.ExactArgs(1),
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
