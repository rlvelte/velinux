package bundle

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	"github.com/rlvelte/velinux/vlx/internal/core/logs"
	"github.com/rlvelte/velinux/vlx/internal/core/su"
	"github.com/rlvelte/velinux/vlx/internal/visuals/notify"
	"github.com/rlvelte/velinux/vlx/internal/visuals/picker"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
	"github.com/rlvelte/velinux/vlx/internal/visuals/progress"
	"github.com/spf13/cobra"
)

// Command returns the cobra command tree for vlx bundle.
func Command() *cobra.Command {
	root := &cobra.Command{
		Use:          "bundle",
		Short:        "Horribly bad bundle installer",
		Long:         "Install/Compile predefined recipes with shell hooks.",
		Aliases:      []string{"bun"},
		Args:         cobra.NoArgs,
		SilenceUsage: !isatty.IsTerminal(os.Stdout.Fd()),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	root.AddCommand(
		&cobra.Command{
			Use:     "list",
			Short:   "List available bundles",
			Long:    "List all available bundles with their segments.",
			Aliases: []string{"l", "ls"},
			Args:    cobra.NoArgs,
			Run:     cmdList,
		},
		&cobra.Command{
			Use:          "install [bundle]",
			Short:        "Install a bundle",
			Long:         "Install a bundle by name or interactively select from a list.",
			Aliases:      []string{"i", "in"},
			Args:         cobra.MaximumNArgs(1),
			SilenceUsage: true,
			Run:          cmdInstall,
		},
	)

	return root
}

// cmdList lists out all bundles and peek at its contents.
func cmdList(cmd *cobra.Command, _ []string) {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)

	dir := fsys.ConfigPath("vlx", "bundles")
	bundles, err := fsys.ListJSON(dir, decodeBundle)
	if err != nil {
		p.Error(err)
		return
	}

	if len(bundles) == 0 {
		p.Success("No bundles found")
		return
	}

	headers := []string{"Name", "Group", "Zypper", "Flatpak", "Hooks"}
	rows := make([][]string, 0, len(bundles))
	for _, b := range bundles {
		hooks := "no"
		if len(b.PreHook) > 0 || len(b.PostHook) > 0 {
			hooks = "yes"
		}

		group := b.Info.Group
		if group == "" {
			group = "-"
		}

		rows = append(rows, []string{
			b.Info.Name,
			group,
			strconv.Itoa(len(b.Zypper)),
			strconv.Itoa(len(b.Flatpak)),
			hooks,
		})
	}

	p.Table(headers, rows)
	return
}

// cmdInstall installs a selected bundle.
func cmdInstall(cmd *cobra.Command, args []string) {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)
	n := cmd.Context().Value(notify.ContextKey).(*notify.Notify)

	dir := fsys.ConfigPath("vlx", "bundles")

	var fileName string
	if len(args) == 0 {
		fn, err := pick(cmd, dir)
		if err != nil {
			p.Error(err)
			return
		}

		fileName = fn
	} else {
		fileName = args[0]
	}

	bundle, err := fsys.GetJSON(dir, fileName, decodeBundle)
	if err != nil {
		p.Error(err)
		return
	}

	steps := []struct {
		enabled bool
		name    string
		run     func() error
	}{
		{len(bundle.Repos) > 0, "Adding repos", func() error { return repo(bundle.Repos) }},
		{len(bundle.PreHook) > 0, "Running pre-install hook", func() error { return hooks(bundle.PreHook) }},
		{len(bundle.Zypper) > 0, "Installing packages", func() error { return zypper(bundle.Zypper) }},
		{len(bundle.Flatpak) > 0, "Installing flatpaks", func() error { return flatpak(bundle.Flatpak) }},
		{len(bundle.PostHook) > 0, "Running post-install hook", func() error { return hooks(bundle.PostHook) }},
	}

	totalSteps := 0
	for _, step := range steps {
		if step.enabled {
			totalSteps++
		}
	}

	prog, progErr := progress.New()
	if progErr == nil && totalSteps > 0 && !isatty.IsTerminal(os.Stdout.Fd()) {
		prog.Start("Installing "+bundle.Info.Name, totalSteps)
		defer prog.Stop()
	}

	for _, step := range steps {
		if !step.enabled {
			continue
		}

		if prog != nil {
			prog.SetLabel(step.name)
		}

		if err := step.run(); err != nil {
			p.Error(err)
			return
		}

		if prog != nil {
			prog.Advance(1)
		}
	}

	_ = n.Send("Successfully installed bundle "+bundle.Info.Name, &notify.Details{
		Title:   "VLX Bundle",
		Urgency: "normal",
	})
}

// pick selects a bundle via an interactive picker.
func pick(cmd *cobra.Command, dir string) (string, error) {
	pkr, err := picker.New()
	if err != nil {
		return "", err
	}

	bundles, err := fsys.ListJSON(dir, decodeBundle)
	if err != nil {
		return "", err
	}

	items := make([]picker.Item, len(bundles))
	lookup := make(map[string]string, len(bundles))
	for i, b := range bundles {
		items[i] = picker.Item{
			Icon:        b.Info.Icon,
			Header:      b.Info.Name,
			Description: b.Info.Description,
			Group:       b.Info.Group,
		}

		lookup[b.Info.Name] = b.filename
	}

	selected, err := pkr.SelectGrouped(cmd.Context(), items)
	if err != nil {
		return "", err
	}

	return lookup[selected.Header], nil
}

// sh executes a cmd through a shell, escalating only when the command.
func sh(str string) error {
	if escalate(str) {
		return su.RunPrivileged("sh", "-c", str)
	}

	cmd := exec.Command("sh", "-c", str)

	cmd.Stdout = logs.Stdout()
	cmd.Stderr = logs.Stderr()

	return cmd.Run()
}

// zypper runs a simple installation with privilege escalation.
func zypper(pkgs []string) error {
	args := append([]string{"zypper", "install", "-y"}, pkgs...)
	return su.RunPrivileged(args...)
}

// flatpak runs a simple installation with privilege escalation.
func flatpak(pkgs []string) error {
	args := append([]string{"flatpak", "install", "-y"}, pkgs...)
	return su.RunPrivileged(args...)
}

// repo adds a repository with privilege escalation.
func repo(repos []Repo) error {
	for _, repo := range repos {
		if err := su.RunPrivileged("zypper", "ar", repo.URL, repo.Alias); err != nil {
			return fmt.Errorf("failed to add repo %q: %w", repo.Alias, err)
		}
	}

	return nil
}

// hooks executes each hook individually with per-command privilege decisions.
func hooks(hooks []string) error {
	for _, hook := range hooks {
		if err := sh(hook); err != nil {
			return fmt.Errorf("hook %q failed: %w", hook, err)
		}
	}

	return nil
}

func escalate(cmd string) bool {
	lower := strings.ToLower(cmd)
	if strings.Contains(lower, "sudo ") || strings.Contains(lower, " sudo ") {
		return true
	}

	if strings.HasPrefix(lower, "su ") || strings.Contains(lower, " su ") {
		return true
	}

	return false
}
