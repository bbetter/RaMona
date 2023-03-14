package utils

import (
	"strings"

	"github.com/mmcdole/gofeed"
)

func FilterByTriggers(items []*gofeed.Item, triggers []string) []*gofeed.Item {

	return Filter(items, func(item *gofeed.Item) bool {
		//equalfold not working!!

		titleLowercase := strings.ToLower(item.Title)
		descLowercase := strings.ToLower(item.Description)

		byTitle := Any(triggers, func(s string) bool {
			return strings.Contains(titleLowercase, s)
		})

		byDescription := Any(triggers, func(s string) bool {
			return strings.Contains(descLowercase, s)
		})

		return byTitle || byDescription
	})
}
