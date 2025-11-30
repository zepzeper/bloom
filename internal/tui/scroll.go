package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Scrolling functions for article view
func scrollDown(m *Model) (*Model, tea.Cmd) {
	if m.CurrentView != "content" {
		return m, nil
	}

	maxScroll := len(m.ArticleLines) - 1
	if m.ScrollOffset < maxScroll {
		m.ScrollOffset++
	}
	return m, nil
}

func scrollUp(m *Model) (*Model, tea.Cmd) {
	if m.CurrentView != "content" {
		return m, nil
	}

	if m.ScrollOffset > 0 {
		m.ScrollOffset--
	}
	return m, nil
}

func scrollPageDown(m *Model) (*Model, tea.Cmd) {
	if m.CurrentView != "content" {
		return m, nil
	}

	// Full screen: height - header (2) - status (1)
	pageSize := max((m.Height-3)/2, 5) // Half screen or at least 5 lines
	maxScroll := len(m.ArticleLines) - 1
	m.ScrollOffset = min(m.ScrollOffset+pageSize, maxScroll)
	return m, nil
}

func scrollPageUp(m *Model) (*Model, tea.Cmd) {
	if m.CurrentView != "content" {
		return m, nil
	}

	// Full screen: height - header (2) - status (1)
	pageSize := max((m.Height-3)/2, 5) // Half screen or at least 5 lines
	m.ScrollOffset = max(m.ScrollOffset-pageSize, 0)
	return m, nil
}

// Helper functions for min/max
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
