package components

import (
	"bloom/internal/feed"
	"bloom/internal/tui/styles"
	"fmt"
	"strings"
)

// RenderArticleList renders the article list view
func RenderArticleList(currentFeed feed.Channel, cursor int, width int) string {
	if len(currentFeed.Item) == 0 {
		return styles.SubtleStyle().Render("No articles in this feed.")
	}

	var items []string
	for i, item := range currentFeed.Item {
		title := item.Title
		if title == "" {
			title = "(Untitled)"
		}

		// Add read indicator
		indicator := "○" // Unread
		if item.Read {
			indicator = "●" // Read
		}

		// Truncate long titles (account for indicator)
		maxTitleWidth := width - 6
		if len(title) > maxTitleWidth {
			title = title[:maxTitleWidth-3] + "..."
		}

		if cursor == i {
			// Selected item - reverse video
			itemText := styles.SelectedStyle().Render(fmt.Sprintf("%s > %s", indicator, title))
			items = append(items, itemText)

			// Show date for selected item
			if item.PubDate != "" {
				dateLine := styles.DateStyle().Render("    " + item.PubDate)
				items = append(items, dateLine)
			}
		} else {
			// Normal item - plain text
			itemText := styles.NormalStyle().Render(fmt.Sprintf("%s   %s", indicator, title))
			items = append(items, itemText)
		}
	}

	return strings.Join(items, "\n")
}

// RenderArticleListStatusBar renders the status bar for article list view
func RenderArticleListStatusBar(feedTitle string, cursor int, articleCount int, width int) string {
	if len(feedTitle) > 30 {
		feedTitle = feedTitle[:27] + "..."
	}
	return styles.RenderStatusBar(
		feedTitle,
		fmt.Sprintf("Article %d/%d", cursor+1, articleCount),
		"↑↓: Navigate  Enter: Read  m: Mark  s: Save  Esc: Back  q: Quit",
		width,
	)
}
