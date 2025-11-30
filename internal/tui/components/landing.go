package components

import (
	"bloom/internal/storage"
	"bloom/internal/tui/styles"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderLanding renders the landing page
func RenderLanding(feedCount, totalArticles, readCount int, config *storage.Config, width, height int) string {
	var sections []string

	// ASCII Art Logo
	logo := `
  ╔══════════════════════════════════════╗
  ║                                      ║
  ║     ██████╗ ██╗      ██████╗  ██╗   ║
  ║     ██╔══██╗██║     ██╔═══██╗██║    ║
  ║     ██████╔╝██║     ██║   ██║██║    ║
  ║     ██╔══██╗██║     ██║   ██║╚═╝    ║
  ║     ██████╔╝███████╗╚██████╔╝██╗    ║
  ║     ╚═════╝ ╚══════╝ ╚═════╝ ╚═╝    ║
  ║                                      ║
  ║         RSS Feed Reader              ║
  ║                                      ║
  ╚══════════════════════════════════════╝
`
	sections = append(sections, styles.ArticleTitleStyle().Render(logo))
	sections = append(sections, "")

	// Stats section
	unreadCount := totalArticles - readCount
	statsBox := renderStatsBox(feedCount, totalArticles, readCount, unreadCount, width)
	sections = append(sections, statsBox)
	sections = append(sections, "")

	// Quick actions
	actionsBox := renderActionsBox(width)
	sections = append(sections, actionsBox)
	sections = append(sections, "")

	// Recent feeds preview
	if config != nil && len(config.Feeds) > 0 {
		feedsPreview := renderFeedsPreview(config.Feeds, width)
		sections = append(sections, feedsPreview)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	
	// Center the content
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

func renderStatsBox(feedCount, totalArticles, readCount, unreadCount int, width int) string {
	boxWidth := min(60, width-4)
	
	title := styles.ArticleTitleStyle().Render("📊 Your Stats")
	
	stats := []string{
		fmt.Sprintf("📡 Feeds:          %d", feedCount),
		fmt.Sprintf("📰 Total Articles: %d", totalArticles),
		fmt.Sprintf("✓  Read:           %d", readCount),
		fmt.Sprintf("○  Unread:         %d", unreadCount),
	}
	
	statsContent := strings.Join(stats, "\n")
	
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(1, 2).
		Width(boxWidth).
		Render(statsContent)
	
	return lipgloss.JoinVertical(lipgloss.Left, title, "", box)
}

func renderActionsBox(width int) string {
	boxWidth := min(60, width-4)
	
	title := styles.ArticleTitleStyle().Render("⚡ Quick Actions")
	
	actions := []string{
		"  f  →  View Feeds",
		"  m  →  Manage Feeds",
		"  s  →  Save State",
		"  q  →  Quit",
	}
	
	actionsContent := strings.Join(actions, "\n")
	
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(1, 2).
		Width(boxWidth).
		Render(actionsContent)
	
	return lipgloss.JoinVertical(lipgloss.Left, title, "", box)
}

func renderFeedsPreview(feeds []storage.FeedConfig, width int) string {
	boxWidth := min(60, width-4)
	
	title := styles.ArticleTitleStyle().Render("📚 Your Feeds")
	
	var feedLines []string
	displayCount := min(5, len(feeds))
	
	for i := 0; i < displayCount; i++ {
		feed := feeds[i]
		// Extract domain from URL for cleaner display
		url := feed.URL
		if len(url) > 40 {
			url = url[:37] + "..."
		}
		
		category := feed.Category
		if category == "" {
			category = "Uncategorized"
		}
		
		feedLines = append(feedLines, fmt.Sprintf("  • %s", url))
		feedLines = append(feedLines, fmt.Sprintf("    %s", styles.SubtleStyle().Render(category)))
	}
	
	if len(feeds) > displayCount {
		feedLines = append(feedLines, "")
		feedLines = append(feedLines, styles.SubtleStyle().Render(fmt.Sprintf("  ... and %d more", len(feeds)-displayCount)))
	}
	
	feedsContent := strings.Join(feedLines, "\n")
	
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(1, 2).
		Width(boxWidth).
		Render(feedsContent)
	
	return lipgloss.JoinVertical(lipgloss.Left, title, "", box)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// RenderLandingStatusBar renders the status bar for the landing page
func RenderLandingStatusBar(width int) string {
	return styles.RenderStatusBar(
		"Welcome",
		"",
		"f: Feeds  m: Manage  s: Save  q: Quit",
		width,
	)
}

