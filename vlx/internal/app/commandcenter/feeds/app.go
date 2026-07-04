package feeds

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	"github.com/rlvelte/velinux/vlx/internal/core/guard"
	"github.com/rlvelte/velinux/vlx/internal/core/printer"
	"github.com/spf13/cobra"
)

// setup validates all requirements for further processing.
func setup(cmd *cobra.Command, _ []string) error {
	if err := errors.Join(guard.Network()); err != nil {
		return err
	}

	cmd.SetContext(context.WithValue(cmd.Context(), printer.ContextKey, printer.New()))
	return nil
}

// Command returns the cobra command tree for vlx commandcenter feeds.
func Command() *cobra.Command {
	root := &cobra.Command{
		Use:               "feeds",
		Short:             "Horribly bad RSS/Atom feed reader",
		Aliases:           []string{"feed", "rss"},
		PersistentPreRunE: setup,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	listCmd := &cobra.Command{
		Use:     "list",
		Short:   "List available feed sources",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE:    cmdFeedsList,
	}

	pollCmd := &cobra.Command{
		Use:     "poll [source]",
		Short:   "Poll feed(s) for latest items",
		Long:    "Poll all feeds or a specific source. Run `vlx cc feeds list` to see available sources.",
		Aliases: []string{"fetch", "get"},
		Args:    cobra.MaximumNArgs(1),
		RunE:    cmdFeedsPoll,
	}
	pollCmd.Flags().BoolP("json", "j", false, "output as JSON")

	root.AddCommand(listCmd)
	root.AddCommand(pollCmd)

	return root
}

func cmdFeedsList(cmd *cobra.Command, args []string) error {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)

	dir := fsys.ConfigPath("vlx", "commandcenter")
	store := fsys.NewStore(dir, decodeFeedsConfig, ".json")
	sources, _ := store.Get("config")

	rows := make([][]string, 0, len(sources))
	for _, s := range sources {
		rows = append(rows, []string{s.Name, s.URL})
	}

	p.Table([]string{"Source", "URL"}, rows)
	return nil
}

func cmdFeedsPoll(cmd *cobra.Command, args []string) error {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)
	jsonFlag, _ := cmd.Flags().GetBool("json")

	dir := fsys.ConfigPath("vlx", "commandcenter")
	store := fsys.NewStore(dir, decodeFeedsConfig, ".json")
	sources, _ := store.Get("config")

	if len(args) > 0 {
		name := args[0]
		var found *FeedSource
		for _, s := range sources {
			if s.Name == name {
				found = &s
				break
			}
		}

		if found == nil {
			return fmt.Errorf("unknown feed source: %s", name)
		}

		sources = []FeedSource{*found}
	}

	type result struct {
		feed *Feed
		err  error
	}

	results := make([]result, len(sources))
	var wg sync.WaitGroup

	for i, s := range sources {
		wg.Add(1)
		go func(i int, s FeedSource) {
			defer wg.Done()

			data, err := get(s.URL)
			if err != nil {
				results[i] = result{err: fmt.Errorf("%s: fetch failed: %w", s.Name, err)}
				return
			}

			feed, err := ParseFeed(s.Name, data)
			if err != nil {
				results[i] = result{err: fmt.Errorf("%s: parse failed: %w", s.Name, err)}
				return
			}

			results[i] = result{feed: feed}
		}(i, s)
	}

	wg.Wait()

	if jsonFlag {
		var feeds []*Feed
		for _, r := range results {
			if r.feed != nil {
				feeds = append(feeds, r.feed)
			}
		}
		data, _ := json.MarshalIndent(feeds, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	for _, r := range results {
		if r.err != nil {
			p.Warn(r.err.Error())
			continue
		}

		p.Info(fmt.Sprintf("\n%s — %s", r.feed.SourceName, r.feed.Title))

		rows := make([][]string, 0, len(r.feed.Items))
		for _, item := range r.feed.Items {
			pub := item.Published
			if pub == "" {
				pub = item.Link
			}
			rows = append(rows, []string{pub, item.Title})
		}

		p.Table([]string{"Posted", "Title"}, rows)
	}

	return nil
}

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
