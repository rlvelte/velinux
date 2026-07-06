package feeds

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	"github.com/rlvelte/velinux/vlx/internal/core/guard"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
	"github.com/spf13/cobra"
)

// setup validates all requirements for further processing.
func setup(_ *cobra.Command, _ []string) error {
	return guard.Network()
}

// Command returns the cobra command tree for vlx fetch feeds.
func Command() *cobra.Command {
	root := &cobra.Command{
		Use:     "feeds",
		Short:   "Horribly bad RSS/Atom feed reader",
		Aliases: []string{"feed", "rss"},
		PreRunE: setup,
		RunE:    cmdPoll,
	}

	listCmd := &cobra.Command{
		Use:     "list",
		Short:   "List available feed sources",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE:    cmdList,
	}

	root.AddCommand(listCmd)
	return root
}

func cmdList(cmd *cobra.Command, _ []string) error {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)

	dir := fsys.ConfigPath("vlx", "fetch")
	store := fsys.NewStore(dir, decodeFeedsConfig, ".json")
	sources, _ := store.Get("config")

	rows := make([][]string, 0, len(sources))
	for _, s := range sources {
		rows = append(rows, []string{s.Name, s.URL})
	}

	p.Table([]string{"Source", "URL"}, rows)
	return nil
}

func cmdPoll(cmd *cobra.Command, args []string) error {
	p := cmd.Context().Value(printer.ContextKey).(*printer.Printer)

	dir := fsys.ConfigPath("vlx", "fetch")
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

			feed, err := parseFeed(s.Name, data)
			if err != nil {
				results[i] = result{err: fmt.Errorf("%s: parse failed: %w", s.Name, err)}
				return
			}

			results[i] = result{feed: feed}
		}(i, s)
	}

	wg.Wait()

	for _, r := range results {
		if r.err != nil {
			p.Warn(r.err.Error())
			continue
		}

		p.Print(fmt.Sprintf("\n%s — %s", r.feed.SourceName, r.feed.Title))

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

func parseFeed(name string, data []byte) (*Feed, error) {
	s := strings.TrimSpace(string(data))
	idx := strings.Index(s, "<?xml")
	if idx > 0 {
		s = s[idx:]
	}

	var rss RSS
	if err := xml.Unmarshal([]byte(s), &rss); err == nil && rss.Channel.Title != "" {
		return rssToFeed(name, &rss), nil
	}

	var atom AtomFeed
	if err := xml.Unmarshal([]byte(s), &atom); err == nil && atom.Title != "" {
		return atomToFeed(name, &atom), nil
	}

	return nil, fmt.Errorf("unrecognized feed format from %s", name)
}

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
