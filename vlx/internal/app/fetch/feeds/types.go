package feeds

import (
	"encoding/json"
)

type FeedSource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Feed struct {
	SourceName  string     `json:"source"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Link        string     `json:"link,omitempty"`
	Items       []FeedItem `json:"items"`
}

type FeedItem struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description,omitempty"`
	Published   string `json:"published,omitempty"`
}

func decodeFeedsConfig(_, _ string, data []byte) ([]FeedSource, error) {
	var wrapper struct {
		Feeds []FeedSource `json:"feeds"`
	}

	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}

	return wrapper.Feeds, nil
}
