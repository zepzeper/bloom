package components

import (
	"bloom/internal/feed"
	"bloom/internal/tui/styles"
	"bloom/internal/tui/utils"
	"fmt"
	"strings"
)

func RenderFeedList(feeds []feed.Channel, currentFeed int, width int) string {
	if len(feeds) == 0 {
		return styles.SubtleStyle().Render("No feeds loaded.")
	}

	var items []string
	for i, feed := range feeds {
		title := feed.Title
		if title == "" {
			title = "(Untitled)"
		}

		if currentFeed == i {
			// Selected item - reverse video
			item := styles.SelectedStyle().Render("> " + title)
			items = append(items, item)

			// Show description for selected item
			if feed.Description != "" {
				desc := utils.StripHTML(feed.Description)
				if len(desc) > width-4 {
					desc = desc[:width-7] + "..."
				}
				descLine := styles.DescriptionStyle().Render("  " + desc)
				items = append(items, descLine)
			}
		} else {
			// Normal item - plain text
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
