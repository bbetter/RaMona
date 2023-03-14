package utils

import (
	"time"

	"github.com/mmcdole/gofeed"
	"golang.org/x/net/context"
)

const (
	timeout = 60 // seconds
	feedUrl = "https://feeds.feedburner.com/gov/gnjU"
)

func ParseFeedItems() []*gofeed.Item {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURLWithContext(feedUrl, ctx)
	return feed.Items
}
