package tui

import (
	"bloom/internal/tui/utils"
	
	"github.com/mattn/go-runewidth"
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

	// Get the current line to calculate display width
	currentLine := m.ArticleLines[actualLine]
	lineDisplayWidth := runewidth.StringWidth(utils.StripANSI(currentLine))
	
	var closestLink *utils.Link
	closestDistance := -1
	
	// Check each link to see if cursor is over it
	for i := range m.ArticleLinks {
		link := &m.ArticleLinks[i]
		if link.Line == actualLine {
			// Check if cursor X is within link bounds
			// Use inclusive bounds: cursor can be at Start or End position
			if m.CursorX >= link.Start && m.CursorX <= link.End {
				return link
			}
			// Also check if cursor is very close to the link (within 1 character)
			// This helps with edge cases where position calculation might be slightly off
			if m.CursorX >= link.Start-1 && m.CursorX <= link.End+1 && m.CursorX >= 0 && m.CursorX <= lineDisplayWidth {
				return link
			}
			
			// Track the closest link on this line as a fallback
			linkCenter := (link.Start + link.End) / 2
			distance := m.CursorX - linkCenter
			if distance < 0 {
				distance = -distance
			}
			if closestLink == nil || distance < closestDistance {
				closestLink = link
				closestDistance = distance
			}
		}
	}

	// If no exact match but we found a link on this line and cursor is reasonably close (within 10 chars), use it
	if closestLink != nil && closestDistance <= 10 {
		return closestLink
	}

	return nil
}
