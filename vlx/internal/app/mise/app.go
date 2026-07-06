package mise

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"sync"

	"github.com/rlvelte/velinux/vlx/internal/core/guard"
	"github.com/rlvelte/velinux/vlx/internal/visuals/notify"
	"github.com/rlvelte/velinux/vlx/internal/visuals/picker"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
	"github.com/spf13/cobra"
)

// setup validates all requirements for further processing.
func setup(cmd *cobra.Command, _ []string) error {
	return errors.Join(guard.Network(), guard.Binaries("mise"))
}

func Command() *cobra.Command {
	root := &cobra.Command{
		Use:     "mise",
		Short:   "Horribly bad mise wrapper",
		Long:    "Manage mise language runtime versions.",
		Args:    cobra.NoArgs,
		Aliases: []string{"lang", "rt"},
		PreRunE: setup,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	root.AddCommand(&cobra.Command{
		Use:     "list",
		Short:   "List installed tool versions",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE:    cmdList,
	})

	root.AddCommand(&cobra.Command{
		Use:     "remote [tool]",
		Short:   "List available remote versions for a tool",
		Aliases: []string{"lsr"},
		Args:    cobra.ExactArgs(1),
		RunE:    cmdRemote,
	})

	root.AddCommand(&cobra.Command{
		Use:     "install [tool] [version]",
		Short:   "Install a tool version",
		Aliases: []string{"i", "add"},
		Args:    cobra.RangeArgs(0, 2),
		RunE:    cmdInstall,
	})

	root.AddCommand(&cobra.Command{
		Use:     "use <tool> <version>",
		Short:   "Set global version for a tool",
		Aliases: []string{"global", "default"},
		Args:    cobra.ExactArgs(2),
		RunE:    cmdUse,
	})

	return root
}

func cmdList(cmd *cobra.Command, _ []string) error {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)

	data, err := exec.Command("mise", "ls", "--json").Output()
	if err != nil {
		return fmt.Errorf("mise ls: %w", err)
	}

	var all miseToolVersions
	if err := json.Unmarshal(data, &all); err != nil {
		return fmt.Errorf("parse mise ls: %w", err)
	}

	var versions []ToolVersion
	for tool, vers := range all {
		for _, v := range vers {
			if !v.Installed {
				continue
			}

			versions = append(versions, ToolVersion{
				Tool:    tool,
				Version: v.Version,
				Active:  v.Active,
				Icon:    iconFor(tool),
			})
		}
	}

	sort.Slice(versions, func(i, j int) bool {
		if versions[i].Tool != versions[j].Tool {
			return versions[i].Tool < versions[j].Tool
		}

		return versions[i].Version > versions[j].Version
	})

	headers := []string{"Tool", "Version", "Active"}
	var rows [][]string
	for _, v := range versions {
		active := ""
		if v.Active {
			active = "*"
		}

		rows = append(rows, []string{v.Tool, v.Version, active})
	}

	p.Table(headers, rows)
	return nil
}

func cmdRemote(cmd *cobra.Command, args []string) error {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)

	tool := args[0]
	data, err := exec.Command("mise", "ls-remote", "--json", tool).Output()
	if err != nil {
		return fmt.Errorf("mise ls-remote: %w", err)
	}

	var versions []RemoteVersion
	if err := json.Unmarshal(data, &versions); err != nil {
		for i, v := range versions {
			versions[i] = RemoteVersion{Version: v.Version, CreatedAt: v.CreatedAt}
		}
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].CreatedAt > versions[j].CreatedAt
	})

	var rows [][]string
	for _, v := range versions {
		rows = append(rows, []string{v.Version, v.CreatedAt[0:10]})
	}

	p.Table([]string{"Version", "Created"}, rows)
	return nil
}

func cmdUse(cmd *cobra.Command, args []string) error {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)
	n := cmd.Context().Value(notify.ContextKey).(*notify.Notify)

	tool := args[0]
	version := args[1]

	out, err := exec.Command("mise", "use", "-g", tool+"@"+version).CombinedOutput()
	if err != nil {
		_ = n.Send(fmt.Sprintf("Failed to set %s to %s", tool, version), &notify.Details{
			Title:   "mise",
			Urgency: "critical",
		})

		return fmt.Errorf("mise use: %w\n%s", err, string(out))
	}

	_ = n.Send(fmt.Sprintf("%s set to %s globally", tool, version), &notify.Details{
		Title:   "mise",
		Urgency: "normal",
	})

	p.Print(fmt.Sprintf("%s %s → %s", iconFor(tool), tool, version))
	return nil
}

func cmdInstall(cmd *cobra.Command, args []string) error {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)
	n := cmd.Context().Value(notify.ContextKey).(*notify.Notify)
	pi := cmd.Context().Value(picker.ContextKey).(*picker.Picker)

	if len(args) > 0 {
		tool := args[0]
		target := tool
		if len(args) > 1 {
			target = tool + "@" + args[1]
		}

		return runInstall(cmd.Context(), p, n, tool, target)
	}

	items, err := buildInstallItems()
	if err != nil {
		return fmt.Errorf("building install items: %w", err)
	}

	if len(items) == 0 {
		p.Warn("No tools available for installation")
		return nil
	}

	selected, err := pi.SelectTwoStage(cmd.Context(), items)
	if err != nil {
		return fmt.Errorf("picker: %w", err)
	}

	target := selected.Header
	parts := strings.SplitN(target, "@", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid selection: %s", target)
	}
	tool := parts[0]

	return runInstall(cmd.Context(), p, n, tool, target)
}

func runInstall(ctx context.Context, p *printer.Printer, n *notify.Notify, tool, target string) error {
	_ = n.Send(fmt.Sprintf("Installing %s...", target), &notify.Details{
		Title:   "mise",
		Urgency: "normal",
	})

	out, err := exec.CommandContext(ctx, "mise", "install", target).CombinedOutput()
	if err != nil {
		_ = n.Send(fmt.Sprintf("Failed to install %s", target), &notify.Details{
			Title:   "mise",
			Urgency: "critical",
		})

		return fmt.Errorf("mise install: %w\n%s", err, string(out))
	}

	resolved := resolveTarget(out, tool)
	if resolved == "" {
		resolved = target
	}

	_ = n.Send(fmt.Sprintf("%s installed successfully", resolved), &notify.Details{
		Title:   "mise",
		Urgency: "normal",
	})

	p.Print(fmt.Sprintf("%s %s installed", iconFor(tool), resolved))
	return nil
}

func buildInstallItems() ([]picker.Item, error) {
	regDesc := fetchRegistryDescriptions()

	tools := make([]string, len(toolOrder))
	copy(tools, toolOrder)

	type result struct {
		tool     string
		versions []RemoteVersion
		err      error
	}

	results := make(chan result, len(tools))
	var wg sync.WaitGroup

	for _, tool := range tools {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			vers, err := fetchRemoteVersions(name)
			results <- result{tool: name, versions: vers, err: err}
		}(tool)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	itemsMap := make(map[string]picker.Item, len(tools))
	for _, t := range tools {
		desc := ""
		if d, ok := regDesc[t]; ok {
			desc = d
		}

		itemsMap[t] = picker.Item{
			Header:      t,
			Description: desc,
			Icon:        toolIcons[t],
		}
	}

	for r := range results {
		if r.err != nil || len(r.versions) == 0 {
			continue
		}

		sort.Slice(r.versions, func(i, j int) bool {
			return r.versions[i].CreatedAt > r.versions[j].CreatedAt
		})

		subitems := make([]picker.Item, 0, 50)
		for i, v := range r.versions {
			if i >= 50 {
				break
			}

			desc := ""
			if len(v.CreatedAt) >= 10 {
				desc = v.CreatedAt[0:10]
			}

			subitems = append(subitems, picker.Item{
				Header:      r.tool + "@" + v.Version,
				Description: desc,
				Icon:        iconFor(r.tool),
			})
		}

		entry := itemsMap[r.tool]
		entry.Subitems = subitems
		itemsMap[r.tool] = entry
	}

	items := make([]picker.Item, 0, len(tools))
	for _, t := range tools {
		items = append(items, itemsMap[t])
	}

	return items, nil
}

func fetchRegistryDescriptions() map[string]string {
	data, err := exec.Command("mise", "registry", "--json").Output()
	if err != nil {
		return nil
	}

	var registry []RegistryTool
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil
	}

	descs := make(map[string]string, len(registry))
	for _, t := range registry {
		descs[t.Short] = t.Description
	}

	return descs
}

func fetchRemoteVersions(tool string) ([]RemoteVersion, error) {
	data, err := exec.Command("mise", "ls-remote", "--json", tool).Output()
	if err != nil {
		return nil, err
	}

	var versions []RemoteVersion
	if err := json.Unmarshal(data, &versions); err != nil {
		return nil, err
	}

	return versions, nil
}

func iconFor(tool string) string {
	if idx := strings.Index(tool, "@"); idx != -1 {
		tool = tool[:idx]
	}

	if i, ok := toolIcons[strings.ToLower(tool)]; ok {
		return i
	}

	return "\uF128"
}

func resolveTarget(out []byte, tool string) string {
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, tool+"@") {
			continue
		}

		if fields := strings.Fields(line); len(fields) > 0 {
			return fields[0]
		}
	}

	return ""
}
