package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating feed request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error receiving response from feed: %w", err)
	}
	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error converting response to byte slice: %w", err)
	}
	data := RSSFeed{}
	err = xml.Unmarshal(bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling bytes into struct: %w", err)
	}
	data.Channel.Title = html.UnescapeString(data.Channel.Title)
	data.Channel.Description = html.UnescapeString(data.Channel.Description)
	for _, item := range data.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}
	return &data, nil
}
