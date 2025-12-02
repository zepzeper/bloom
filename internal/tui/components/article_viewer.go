package components

import (
	"bloom/internal/tui/styles"
	"bloom/internal/tui/utils"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
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

		// Add cursor indicator to the current line (preserving ANSI codes)
		if cursorY >= 0 && cursorY < len(visibleLines) {
			line := visibleLines[cursorY]
			// Get display width to validate cursor position
			displayWidth := runewidth.StringWidth(utils.StripANSI(line))
			cursorPos := cursorX
			if cursorPos > displayWidth {
				cursorPos = displayWidth
			}
			// Insert cursor while preserving ANSI color codes
			visibleLines[cursorY] = utils.InsertCursorAtPosition(line, cursorPos)
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
	helpText := "o/O:Open c:Copy"
	if currentLink != nil {
		// Show link URL if over a link
		url := currentLink.URL
		// Calculate available space for right side
		// Reserve space for help text: "[url] o/O:Open c:Copy"
		leftWidth := len(left)
		centerWidth := len(center)
		helpTextLen := len(helpText) + 3 // "[url] " prefix
		availableForURL := width - leftWidth - centerWidth - helpTextLen - 6 // 6 for spacing
		
		if availableForURL < 5 {
			// Very narrow terminal, just show help
			right = helpText
		} else {
			if len(url) > availableForURL {
				url = url[:availableForURL-3] + "..."
			}
			right = fmt.Sprintf("[%s] %s", url, helpText)
		}
	} else {
		// Show vim navigation help
		right = "j/k:Scroll w/b:Word 0/$:Line Esc:Back q:Quit"
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
		// Not enough space, prioritize showing help text
		// Truncate URL if needed, but keep help text visible
		if currentLink != nil && len(right) > width-leftWidth-centerWidth-2 {
			// Recalculate right side with better truncation
			maxURLLen := width - leftWidth - centerWidth - len(helpText) - 8
			if maxURLLen < 3 {
				right = helpText // Just show help if too narrow
			} else {
				url := currentLink.URL
				if len(url) > maxURLLen {
					url = url[:maxURLLen-3] + "..."
				}
				right = fmt.Sprintf("[%s] %s", url, helpText)
			}
		}
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
