package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/rlvelte/velinux/vlx/internal/app/bundle"
	"github.com/rlvelte/velinux/vlx/internal/app/fetch"
	"github.com/rlvelte/velinux/vlx/internal/app/launcher"
	"github.com/rlvelte/velinux/vlx/internal/app/themes"
	"github.com/rlvelte/velinux/vlx/internal/core/logs"
	"github.com/rlvelte/velinux/vlx/internal/visuals/notify"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
	"github.com/spf13/cobra"
)

var session *logs.Session

func setup(cmd *cobra.Command, _ []string) {
	call := strings.ReplaceAll(cmd.CommandPath(), " ", "-")
	if s, err := logs.Open(call); err == nil {
		session = s
		slog.SetDefault(s.Logger())
	}

	p := printer.New()
	if jsonFlag, _ := cmd.Flags().GetBool("json"); jsonFlag {
		p.ForceJSON()
	}
	ctx := context.WithValue(cmd.Context(), printer.ContextKey, p)

	ctx = context.WithValue(ctx, notify.ContextKey, notify.New())
	cmd.SetContext(ctx)
}

func main() {
	root := &cobra.Command{
		Use:              "vlx",
		Short:            "Horribly bad utilities",
		Long:             "VeLinux centered command utility application.",
		SilenceUsage:     !isatty.IsTerminal(os.Stdout.Fd()),
		PersistentPreRun: setup,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	root.PersistentFlags().BoolP("json", "j", false, "output as JSON")

	root.AddCommand(
		completions(),
		themes.Command(),
		launcher.Command(),
		bundle.Command(),
		fetch.Command(),
	)

	_ = root.Execute()
	_ = session.Close()
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
