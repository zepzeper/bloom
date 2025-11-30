package components

import (
	"bloom/internal/feed"
	"bloom/internal/storage"
	"bloom/internal/tui/styles"
	"bloom/internal/tui/utils"
	"fmt"
	"strings"
)

func RenderFeedList(configFeeds []storage.FeedConfig, loadedFeeds []feed.Channel, currentFeed int, width int) string {
	if len(configFeeds) == 0 {
		return styles.SubtleStyle().Render("No feeds configured. Press 'm' to manage feeds.")
	}

		var items []string
	for i, feedConfig := range configFeeds {
		// Find matching loaded feed by FeedURL (the URL we used to fetch it)
		var loadedFeed *feed.Channel
		for j := range loadedFeeds {
			if loadedFeeds[j].FeedURL == feedConfig.URL {
				loadedFeed = &loadedFeeds[j]
				break
			}
		}

		var title string
		var description string
		var isLoaded bool

		if loadedFeed != nil {
			// Feed is loaded
			isLoaded = true
			title = loadedFeed.Title
			if title == "" {
				title = "(Untitled)"
			}
			description = loadedFeed.Description
		} else {
			// Feed not loaded yet or failed to load
			isLoaded = false
			// Use URL as title if feed hasn't loaded
			title = feedConfig.URL
			if len(title) > width-20 {
				title = title[:width-23] + "..."
			}
			title = title + " (Loading...)"
		}

		if currentFeed == i {
			// Selected item
			item := styles.SelectedStyle().Render("> " + title)
			items = append(items, item)

			// Show description for selected item if loaded
			if isLoaded && description != "" {
				desc := utils.StripHTML(description)
				if len(desc) > width-4 {
					desc = desc[:width-7] + "..."
				}
				descLine := styles.DescriptionStyle().Render("  " + desc)
				items = append(items, descLine)
			}
		} else {
			// Normal item
			item := styles.NormalStyle().Render("  " + title)
			items = append(items, item)
		}
	}

	return strings.Join(items, "\n")
}

func RenderFeedStatusBar(feedCount int, width int) string {
	return styles.RenderStatusBar(
		"Feeds",
		fmt.Sprintf("%d feed(s)", feedCount),
		"↑↓: Navigate  Enter: Open  f: Manage  Esc: Home  q: Quit",
		width,
	)
}
