package tui

import (
	"bloom/internal/tui/components"
	"bloom/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

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
		
		return components.RenderLanding(
			len(m.Feeds),
			totalArticles,
			readCount,
			m.Config,
			width,
			m.Height,
		) + "\n" + components.RenderLandingStatusBar(width)
	case "feed":
		content = components.RenderFeedList(m.Feeds, m.CurrentFeed, width)
		status = components.RenderFeedStatusBar(len(m.Feeds), width)
		return lipgloss.JoinVertical(lipgloss.Left, content, status)
	case "articles":
		if m.CurrentFeed < len(m.Feeds) {
			feed := m.Feeds[m.CurrentFeed]
			content = components.RenderArticleList(feed, m.Cursor, width)
			status = components.RenderArticleListStatusBar(feed.Title, m.Cursor, len(feed.Item), width)
			return lipgloss.JoinVertical(lipgloss.Left, content, status)
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
