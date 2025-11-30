package handlers

import (
	"bloom/internal/tui"
	"bloom/internal/tui/utils"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

func HandleFeedLoad(m *tui.Model, msg tui.FeedLoadMsg) (*tui.Model, tea.Cmd) {
	m.Loading = false

	if msg.Err != nil {
		m.Err = msg.Err
		return m, nil
	}

	if msg.Channel != nil {
		m.Feeds = append(m.Feeds, *msg.Channel)
		m.Err = nil
	}

	return m, nil
}

func HandleArticleLoad(m *tui.Model, msg tui.ArticleLoadMsg) (*tui.Model, tea.Cmd) {
	m.Loading = false

	if msg.Err != nil {
		m.Err = msg.Err
		return m, nil
	}

	m.CurrentArticle = msg.Article
	m.ArticleContent = msg.Article.Content

	// Render markdown with glamour and parse into lines for scrolling
	renderedContent, err := renderMarkdownForScrolling(msg.Article.Content, m.Width-8)
	if err != nil {
		// Fallback to raw content if rendering fails
		m.ArticleLines = strings.Split(msg.Article.Content, "\n")
		m.ArticleLinks = utils.ParseLinksFromMarkdown(msg.Article.Content)
	} else {
		m.ArticleLines = strings.Split(renderedContent, "\n")
		// Parse links from rendered content (strip ANSI codes first for accurate positioning)
		m.ArticleLinks = utils.ParseLinksFromRenderedContent(renderedContent)
	}

	m.ScrollOffset = 0 // Reset scroll to top
	m.CursorX = 0      // Reset cursor position
	m.CursorY = 0      // Reset cursor position
	m.CurrentView = "content"
	m.Err = nil

	return m, nil
}

// renderMarkdownForScrolling renders markdown with glamour for line counting
func renderMarkdownForScrolling(markdownContent string, width int) (string, error) {
	if markdownContent == "" {
		return "", nil
	}

	// Render with glamour
	rendered, err := glamour.RenderWithEnvironmentConfig(markdownContent)
	if err != nil {
		return markdownContent, err
	}

	return rendered, nil
}
