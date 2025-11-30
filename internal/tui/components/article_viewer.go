package components

import (
	"bloom/internal/tui/styles"
	"bloom/internal/tui/utils"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func RenderArticleFullScreen(title string, articleLines []string, articleLinks []utils.Link, scrollOffset, cursorX, cursorY, width, height int) string {
	if len(articleLines) == 0 {
		return styles.SubtleStyle().Render("No content available.")
	}

	var parts []string

	headerHeight := 1
	if title != "" {
		displayTitle := title
		if len(displayTitle) > width {
			displayTitle = displayTitle[:width-3] + "..."
		}
		parts = append(parts, styles.ArticleTitleStyle().Render(displayTitle))
		headerHeight = 2
	}

	// (full screen minus header and status)
	statusHeight := 1
	contentHeight := height - headerHeight - statusHeight
	if contentHeight < 1 {
		contentHeight = height - 2
	}

	// Get visible lines
	start := scrollOffset
	end := start + contentHeight
	if end > len(articleLines) {
		end = len(articleLines)
	}

	if start < len(articleLines) && len(articleLines) > 0 {
		visibleLines := make([]string, len(articleLines[start:end]))
		copy(visibleLines, articleLines[start:end])

		// Highlight link under cursor if any
		currentLink := findLinkAtPosition(articleLines, articleLinks, scrollOffset, cursorX, cursorY)
		if currentLink != nil && currentLink.Line >= start && currentLink.Line < end {
			lineIdx := currentLink.Line - start
			if lineIdx >= 0 && lineIdx < len(visibleLines) {
				line := visibleLines[lineIdx]
				urlInLine := utils.FindURLInRenderedLine(line, currentLink.URL)
				if urlInLine != "" {
					highlightedURL := styles.LinkStyle().Render(currentLink.URL)
					visibleLines[lineIdx] = strings.Replace(line, urlInLine, highlightedURL, 1)
				}
			}
		}

		// Add cursor indicator to the current line
		if cursorY >= 0 && cursorY < len(visibleLines) {
			line := visibleLines[cursorY]
			// Strip ANSI codes to get actual display width
			cleanLine := utils.StripANSI(line)
			displayWidth := len([]rune(cleanLine))
			
			// Add cursor at position (or at end if position is beyond line)
			cursorPos := cursorX
			if cursorPos > displayWidth {
				cursorPos = displayWidth
			}
			
			// Insert cursor marker (inverse video space)
			if cursorPos < displayWidth {
				runes := []rune(cleanLine)
				before := string(runes[:cursorPos])
				after := string(runes[cursorPos:])
				cursorChar := lipgloss.NewStyle().Reverse(true).Render(string(runes[cursorPos]))
				visibleLines[cursorY] = before + cursorChar + after[1:]
			} else if cursorPos == displayWidth {
				// Cursor at end of line
				visibleLines[cursorY] = cleanLine + lipgloss.NewStyle().Reverse(true).Render(" ")
			}
		}

		// Pad to fill screen
		for len(visibleLines) < contentHeight {
			visibleLines = append(visibleLines, "")
		}

		parts = append(parts, strings.Join(visibleLines, "\n"))
	}

	// Add minimal status line (man-page style)
	currentLink := findLinkAtPosition(articleLines, articleLinks, scrollOffset, cursorX, cursorY)
	status := RenderManPageStatusBar(title, scrollOffset, len(articleLines), currentLink, width)

	content := strings.Join(parts, "\n")
	return lipgloss.JoinVertical(lipgloss.Left, content, status)
}

func RenderManPageStatusBar(title string, scrollOffset, totalLines int, currentLink *utils.Link, width int) string {
	displayTitle := title
	if len(displayTitle) > 25 {
		displayTitle = displayTitle[:22] + "..."
	}

	// Build status line (man-page style)
	var left, center, right string

	// Left: Article name
	left = displayTitle

	// Center: Line position
	if totalLines > 0 {
		center = fmt.Sprintf("line %d/%d", scrollOffset+1, totalLines)
	}

	// Right: Link info or help
	if currentLink != nil {
		// Show link URL if over a link
		url := currentLink.URL
		if len(url) > 25 {
			url = url[:22] + "..."
		}
		right = fmt.Sprintf("[%s] o:Open c:Copy", url)
	} else {
		// Show minimal help
		right = "j/k:Scroll Esc:Back q:Quit"
	}

	// Format like man page with proper spacing
	// Calculate spacing
	leftWidth := len(left)
	centerWidth := len(center)
	rightWidth := len(right)

	// Center the center part
	availableWidth := width - leftWidth - rightWidth - 4 // 4 for spacing
	if availableWidth > centerWidth {
		leftSpacer := (availableWidth - centerWidth) / 2
		rightSpacer := availableWidth - centerWidth - leftSpacer
		statusLine := left + strings.Repeat(" ", leftSpacer+1) + center + strings.Repeat(" ", rightSpacer+1) + right
		// Pad to exact width
		if len(statusLine) < width {
			statusLine = statusLine + strings.Repeat(" ", width-len(statusLine))
		} else if len(statusLine) > width {
			statusLine = statusLine[:width]
		}
		return styles.StatusStyle().Width(width).Render(statusLine)
	} else {
		// Not enough space, just concatenate
		statusLine := left + " " + center + " " + right
		if len(statusLine) > width {
			statusLine = statusLine[:width]
		}
		return styles.StatusStyle().Width(width).Render(statusLine)
	}
}

// findLinkAtPosition finds a link at the given cursor position
func findLinkAtPosition(articleLines []string, articleLinks []utils.Link, scrollOffset, cursorX, cursorY int) *utils.Link {
	// Calculate actual line in ArticleLines (accounting for scroll)
	actualLine := scrollOffset + cursorY
	if actualLine < 0 || actualLine >= len(articleLines) {
		return nil
	}

	// Check each link to see if cursor is over it
	for i := range articleLinks {
		link := &articleLinks[i]
		if link.Line == actualLine {
			// Check if cursor X is within link bounds
			if cursorX >= link.Start && cursorX <= link.End {
				return link
			}
		}
	}

	return nil
}
