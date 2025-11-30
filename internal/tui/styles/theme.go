package styles

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var (
	// Classic terminal colors - simple and minimal
	bgColor     = lipgloss.Color("0") // Black
	fgColor     = lipgloss.Color("7") // White/Light gray
	selectedBg  = lipgloss.Color("7") // White for selection
	selectedFg  = lipgloss.Color("0") // Black text on selection
	subtleColor = lipgloss.Color("8") // Dark gray
	borderColor = lipgloss.Color("8") // Dark gray for borders
	linkColor   = lipgloss.Color("4") // Blue for links

	// Normal text style - plain white on black
	normalStyle = lipgloss.NewStyle().
			Foreground(fgColor).
			Background(bgColor)

	// Selected item - reverse video (classic terminal style)
	selectedStyle = lipgloss.NewStyle().
			Foreground(selectedFg).
			Background(selectedBg).
			Reverse(true)

	// Cursor indicator - simple ">"
	cursorStyle = lipgloss.NewStyle().
			Foreground(fgColor).
			SetString(">")

	// Subtle text
	subtleStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	// Title style - just bold, no fancy colors
	titleStyle = lipgloss.NewStyle().
			Foreground(fgColor).
			Bold(true)

	// Link style - underline for links
	linkStyle = lipgloss.NewStyle().
			Foreground(linkColor).
			Underline(true)

	// Status bar - simple line
	statusStyle = lipgloss.NewStyle().
			Foreground(subtleColor).
			Background(bgColor)

	// Error style - just red
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1"))

	// Loading style - subtle
	loadingStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	// Article title - bold white
	articleTitleStyle = lipgloss.NewStyle().
			Foreground(fgColor).
			Bold(true)

	// Article meta - subtle gray
	articleMetaStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	// Feed description - subtle
	descriptionStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	// Date style - subtle
	dateStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	// Cursor style for article viewer
	cursorBlockStyle = lipgloss.NewStyle().
				Foreground(selectedFg).
				Background(selectedBg).
				Reverse(true)
)

// Exported style getters
func NormalStyle() lipgloss.Style {
	return normalStyle
}

func SelectedStyle() lipgloss.Style {
	return selectedStyle
}

func SubtleStyle() lipgloss.Style {
	return subtleStyle
}

func ErrorStyle() lipgloss.Style {
	return errorStyle
}

func LoadingStyle() lipgloss.Style {
	return loadingStyle
}

func ArticleTitleStyle() lipgloss.Style {
	return articleTitleStyle
}

func DescriptionStyle() lipgloss.Style {
	return descriptionStyle
}

func DateStyle() lipgloss.Style {
	return dateStyle
}

func LinkStyle() lipgloss.Style {
	return linkStyle
}

func StatusStyle() lipgloss.Style {
	return statusStyle
}

func CursorStyle() lipgloss.Style {
	return cursorBlockStyle
}

// RenderStatusBar creates a simple status bar
func RenderStatusBar(view string, position string, help string, width int) string {
	left := subtleStyle.Render(view)
	center := subtleStyle.Render(position)
	right := subtleStyle.Render(help)

	// Simple separator line
	separator := strings.Repeat("─", width)

	// Calculate spacing
	availableWidth := max(width-lipgloss.Width(left)-lipgloss.Width(right)-4, 0)
	spacer := strings.Repeat(" ", availableWidth/2)
	centered := spacer + center + spacer

	content := lipgloss.JoinHorizontal(lipgloss.Left, left, centered, right)
	return statusStyle.Width(width).Render(separator + "\n" + content)
}

