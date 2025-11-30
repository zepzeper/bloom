package components

import (
	"bloom/internal/tui/styles"
	"fmt"
	"strings"
)

// RenderCategoryList renders a list of categories with feed counts
func RenderCategoryList(categories []string, current int, width int) string {
	if len(categories) == 0 {
		return styles.SubtleStyle().Render("No categories available.")
	}

	var items []string
	for i, category := range categories {
		displayName := category
		if displayName == "" {
			displayName = "All"
		}

		// Truncate long category names
		maxWidth := width - 4
		if len(displayName) > maxWidth {
			displayName = displayName[:maxWidth-3] + "..."
		}

		if current == i {
			// Selected category
			itemText := styles.SelectedStyle().Render("> " + displayName)
			items = append(items, itemText)
		} else {
			// Normal category
			itemText := styles.NormalStyle().Render("  " + displayName)
			items = append(items, itemText)
		}
	}

	return strings.Join(items, "\n")
}

// RenderCategoryStatusBar renders the status bar for category view
func RenderCategoryStatusBar(categoryCount int, width int) string {
	return styles.RenderStatusBar(
		"Categories",
		fmt.Sprintf("%d categories", categoryCount),
		"↑↓: Navigate  Enter: Select  Esc: Back  q: Quit",
		width,
	)
}
