package components

import (
	"bloom/internal/storage"
	"bloom/internal/tui/styles"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderFeedManager renders the feed management view
func RenderFeedManager(feeds []storage.FeedConfig, cursor int, editing bool, editField string, editValue string, width int) string {
	if len(feeds) == 0 {
		return styles.SubtleStyle().Render("No feeds configured. Press 'a' to add a feed.")
	}

	var items []string
	
	// Header
	header := styles.ArticleTitleStyle().Render("Feed Management")
	items = append(items, header, "")

	for i, feed := range feeds {
		// Check if this feed is being edited
		isEditing := editing && cursor == i

		var feedDisplay string
		if isEditing {
			// Show edit form
			feedDisplay = renderEditForm(feed, editField, editValue, width)
		} else {
			// Show normal feed info
			feedDisplay = renderFeedInfo(feed, cursor == i, width)
		}

		items = append(items, feedDisplay)
		
		// Add spacing between feeds
		if i < len(feeds)-1 {
			items = append(items, "")
		}
	}

	return strings.Join(items, "\n")
}

func renderFeedInfo(feed storage.FeedConfig, selected bool, width int) string {
	url := feed.URL
	if len(url) > width-10 {
		url = url[:width-13] + "..."
	}

	category := feed.Category
	if category == "" {
		category = "Uncategorized"
	}

	tags := strings.Join(feed.Tags, ", ")
	if tags == "" {
		tags = "No tags"
	}

	var lines []string
	if selected {
		lines = append(lines, styles.SelectedStyle().Render(fmt.Sprintf("> %s", url)))
		lines = append(lines, styles.SubtleStyle().Render(fmt.Sprintf("  Category: %s", category)))
		lines = append(lines, styles.SubtleStyle().Render(fmt.Sprintf("  Tags: %s", tags)))
	} else {
		lines = append(lines, styles.NormalStyle().Render(fmt.Sprintf("  %s", url)))
		lines = append(lines, styles.SubtleStyle().Render(fmt.Sprintf("  Category: %s | Tags: %s", category, tags)))
	}

	return strings.Join(lines, "\n")
}

func renderEditForm(feed storage.FeedConfig, editField string, editValue string, width int) string {
	var lines []string

	// URL field
	urlLabel := "URL: "
	urlValue := feed.URL
	if editField == "url" {
		urlValue = editValue + "█" // Show cursor
		urlLabel = "> " + urlLabel
	}
	lines = append(lines, styles.ArticleTitleStyle().Render(urlLabel)+urlValue)

	// Category field
	categoryLabel := "Category: "
	categoryValue := feed.Category
	if categoryValue == "" {
		categoryValue = "(empty)"
	}
	if editField == "category" {
		categoryValue = editValue + "█"
		categoryLabel = "> " + categoryLabel
	}
	lines = append(lines, styles.NormalStyle().Render(categoryLabel)+categoryValue)

	// Tags field
	tagsLabel := "Tags: "
	tagsValue := strings.Join(feed.Tags, ", ")
	if tagsValue == "" {
		tagsValue = "(empty)"
	}
	if editField == "tags" {
		tagsValue = editValue + "█"
		tagsLabel = "> " + tagsLabel
	}
	lines = append(lines, styles.NormalStyle().Render(tagsLabel)+tagsValue)

	// Instructions
	lines = append(lines, "")
	lines = append(lines, styles.SubtleStyle().Render("Tab: Next field | Ctrl+V: Paste | Enter: Save | Esc: Cancel"))

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// RenderFeedManagerStatusBar renders the status bar for feed management view
func RenderFeedManagerStatusBar(feedCount int, width int) string {
	return styles.RenderStatusBar(
		"Feed Manager",
		fmt.Sprintf("%d feeds", feedCount),
		"a: Add  e: Edit  d: Delete  r: Reload  Esc: Home  q: Quit",
		width,
	)
}

// RenderAddFeedForm renders the form for adding a new feed
func RenderAddFeedForm(url, category, tags string, currentField string, width int) string {
	var lines []string

	// Title
	title := styles.ArticleTitleStyle().Render("Add New Feed")
	lines = append(lines, title, "")

	// URL field
	urlLabel := "URL: "
	urlValue := url
	if currentField == "url" {
		urlValue = url + "█"
		urlLabel = "> " + urlLabel
	}
	if urlValue == "" || urlValue == "█" {
		urlValue = "(required)"
	}
	lines = append(lines, styles.NormalStyle().Render(urlLabel)+urlValue)

	// Category field
	categoryLabel := "Category: "
	categoryValue := category
	if currentField == "category" {
		categoryValue = category + "█"
		categoryLabel = "> " + categoryLabel
	}
	if categoryValue == "" || categoryValue == "█" {
		categoryValue = "(optional)"
	}
	lines = append(lines, styles.NormalStyle().Render(categoryLabel)+categoryValue)

	// Tags field
	tagsLabel := "Tags (comma-separated): "
	tagsValue := tags
	if currentField == "tags" {
		tagsValue = tags + "█"
		tagsLabel = "> " + tagsLabel
	}
	if tagsValue == "" || tagsValue == "█" {
		tagsValue = "(optional)"
	}
	lines = append(lines, styles.NormalStyle().Render(tagsLabel)+tagsValue)

	// Instructions
	lines = append(lines, "")
	lines = append(lines, styles.SubtleStyle().Render("Tab: Next field | Ctrl+V: Paste | Enter: Save | Esc: Cancel"))

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

