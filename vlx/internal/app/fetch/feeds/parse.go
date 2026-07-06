package feeds

import (
	"encoding/xml"
	"strings"
	"time"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func rssToFeed(name string, rss *RSS) *Feed {
	feed := &Feed{
		SourceName:  name,
		Title:       rss.Channel.Title,
		Description: rss.Channel.Description,
		Link:        rss.Channel.Link,
	}

	for _, item := range rss.Channel.Items {
		pub := item.PubDate
		if t, err := time.Parse(time.RFC1123Z, pub); err == nil {
			pub = t.Format("2006-01-02 15:04")
		} else if t, err := time.Parse(time.RFC1123, pub); err == nil {
			pub = t.Format("2006-01-02 15:04")
		}

		feed.Items = append(feed.Items, FeedItem{
			Title:       item.Title,
			Link:        item.Link,
			Description: cleanDesc(item.Description),
			Published:   pub,
		})
	}

	return feed
}

type AtomFeed struct {
	XMLName xml.Name    `xml:"http://www.w3.org/2005/Atom feed"`
	Title   string      `xml:"title"`
	Link    []AtomLink  `xml:"link"`
	Entries []AtomEntry `xml:"entry"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
}

type AtomEntry struct {
	Title     string      `xml:"title"`
	Link      []AtomLink  `xml:"link"`
	Published string      `xml:"published"`
	Updated   string      `xml:"updated"`
	Summary   string      `xml:"summary"`
	Content   AtomContent `xml:"content"`
}

type AtomContent struct {
	Text string `xml:",chardata"`
}

func atomToFeed(name string, atom *AtomFeed) *Feed {
	feed := &Feed{
		SourceName: name,
		Title:      atom.Title,
	}

	for _, l := range atom.Link {
		if l.Href != "" {
			feed.Link = l.Href
			break
		}
	}

	for _, entry := range atom.Entries {
		pub := entry.Published
		if pub == "" {
			pub = entry.Updated
		}
		if t, err := time.Parse(time.RFC3339, pub); err == nil {
			pub = t.Format("2006-01-02 15:04")
		} else if t, err := time.Parse(time.RFC3339Nano, pub); err == nil {
			pub = t.Format("2006-01-02 15:04")
		}

		desc := entry.Summary
		if desc == "" {
			desc = entry.Content.Text
		}

		link := ""
		for _, l := range entry.Link {
			if l.Href != "" {
				link = l.Href
				break
			}
		}

		feed.Items = append(feed.Items, FeedItem{
			Title:       entry.Title,
			Link:        link,
			Description: cleanDesc(desc),
			Published:   pub,
		})
	}

	return feed
}

func cleanDesc(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 200 {
		s = s[:200] + "..."
	}
	return s
}
