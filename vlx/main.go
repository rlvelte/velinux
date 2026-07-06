package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/rlvelte/velinux/vlx/internal/app/bundle"
	"github.com/rlvelte/velinux/vlx/internal/app/fetch"
	"github.com/rlvelte/velinux/vlx/internal/app/launcher"
	"github.com/rlvelte/velinux/vlx/internal/app/mise"
	"github.com/rlvelte/velinux/vlx/internal/app/themes"
	"github.com/rlvelte/velinux/vlx/internal/visuals/notify"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
	"github.com/spf13/cobra"
)

func setup(cmd *cobra.Command, args []string) error {
	p := printer.New()
	if jsonFlag, _ := cmd.Flags().GetBool("json"); jsonFlag {
		p = p.ForceJSON()
	}

	cmd.SetContext(context.WithValue(cmd.Context(), printer.ContextKey, p))
	cmd.SetContext(context.WithValue(cmd.Context(), notify.ContextKey, notify.New()))
}

func main() {
	root := &cobra.Command{
		Use:               "vlx",
		Short:             "Horribly bad utilities",
		Long:              "VeLinux centered command utility application.",
		PersistentPreRunE: setup,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	root.PersistentFlags().BoolP("json", "j", false, "output as JSON")

	root.AddCommand(
		completions(),
		themes.Command(),
		launcher.Command(),
		mise.Command(),
		bundle.Command(),
		fetch.Command(),
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
