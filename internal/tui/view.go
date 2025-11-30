package tui

import (
	"bloom/internal/feed"
	"bloom/internal/tui/components"
	"bloom/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// getLoadedFeedForConfigIndex returns the loaded feed for the given config feed index,
// or nil if the feed hasn't loaded yet or doesn't exist
func (m Model) getLoadedFeedForConfigIndex(configIndex int) *feed.Channel {
	if m.Config == nil || configIndex < 0 || configIndex >= len(m.Config.Feeds) {
		return nil
	}
	
	feedConfig := m.Config.Feeds[configIndex]
	// Find matching loaded feed by FeedURL
	for i := range m.Feeds {
		if m.Feeds[i].FeedURL == feedConfig.URL {
			return &m.Feeds[i]
		}
	}
	return nil
}

// View is the main view dispatcher (bubbletea interface)
func (m Model) View() string {
	width := m.Width
	if width < 40 {
		width = 80 // Default width
	}

	// Error view
	if m.Err != nil {
		return components.RenderErrorView(m.Err, width)
	}

	// Loading view
	if m.Loading {
		return components.RenderLoadingView(width)
	}

	// Main content views
	var content string
	var status string

	switch m.CurrentView {
	case "landing":
		// Calculate stats
		totalArticles := 0
		readCount := 0
		for _, feed := range m.Feeds {
			totalArticles += len(feed.Item)
			for _, item := range feed.Item {
				if item.Read {
					readCount++
				}
			}
		}

		// Use Config.Feeds count to match what's displayed in the feed list
		feedCount := len(m.Config.Feeds)
		if m.Config == nil {
			feedCount = 0
		}

		return components.RenderLanding(
			feedCount,
			totalArticles,
			readCount,
			m.Config,
			width,
			m.Height,
		) + "\n" + components.RenderLandingStatusBar(width)
	case "feed":
		feedCount := len(m.Config.Feeds)
		if m.Config == nil {
			feedCount = 0
		}
		content = components.RenderFeedList(m.Config.Feeds, m.Feeds, m.CurrentFeed, width)
		status = components.RenderFeedStatusBar(feedCount, width)
		return lipgloss.JoinVertical(lipgloss.Left, content, status)
	case "articles":
		loadedFeed := m.getLoadedFeedForConfigIndex(m.CurrentFeed)
		if loadedFeed != nil {
			content = components.RenderArticleList(*loadedFeed, m.Cursor, width)
			status = components.RenderArticleListStatusBar(loadedFeed.Title, m.Cursor, len(loadedFeed.Item), width)
			return lipgloss.JoinVertical(lipgloss.Left, content, status)
		}
		// Feed not loaded yet
		if m.Config != nil && m.CurrentFeed < len(m.Config.Feeds) {
			return styles.SubtleStyle().Render("Feed is loading...") + "\n" + styles.RenderStatusBar("Articles", "Loading...", "Esc: Back", width)
		}
		return "No feed selected"
	case "content":
		// Full-screen man-page style view
		return components.RenderArticleFullScreen(
			m.CurrentArticle.Title,
			m.ArticleLines,
			m.ArticleLinks,
			m.ScrollOffset,
			m.CursorX,
			m.CursorY,
			width,
			m.Height,
		)
	case "manage":
		// Feed management view
		if m.AddingFeed {
			content = components.RenderAddFeedForm(
				m.AddFeedURL,
				m.AddFeedCat,
				m.AddFeedTags,
				m.AddFeedField,
				width,
			)
			status = styles.RenderStatusBar("Add Feed", "", "Tab: Next | Enter: Save | Esc: Cancel", width)
		} else {
			content = components.RenderFeedManager(
				m.Config.Feeds,
				m.Cursor,
				m.EditingFeed,
				m.EditField,
				m.EditValue,
				width,
			)
			status = components.RenderFeedManagerStatusBar(len(m.Config.Feeds), width)
		}
		return lipgloss.JoinVertical(lipgloss.Left, content, status)
	default:
		content = "Unknown view"
		status = styles.RenderStatusBar("Unknown", "", "Press q to quit", width)
		return lipgloss.JoinVertical(lipgloss.Left, content, status)
	}
}
