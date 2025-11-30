package tui

import (
	"bloom/internal/tui/utils"
)

// findLinkAtPosition finds a link at the given cursor position
func findLinkAtPosition(m *Model) *utils.Link {
	if m.CurrentView != "content" {
		return nil
	}

	// Calculate actual line in ArticleLines (accounting for scroll)
	actualLine := m.ScrollOffset + m.CursorY
	if actualLine < 0 || actualLine >= len(m.ArticleLines) {
		return nil
	}

	// Check each link to see if cursor is over it
	for i := range m.ArticleLinks {
		link := &m.ArticleLinks[i]
		if link.Line == actualLine {
			// Check if cursor X is within link bounds
			if m.CursorX >= link.Start && m.CursorX <= link.End {
				return link
			}
		}
	}

	return nil
}
