package tui

import (
	"bloom/internal/storage"
	"bloom/internal/tui/utils"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

// Message handlers
func handleFeedLoad(m *Model, msg FeedLoadMsg) (*Model, tea.Cmd) {
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

func handleArticleLoad(m *Model, msg ArticleLoadMsg) (*Model, tea.Cmd) {
	m.Loading = false

	if msg.Err != nil {
		m.Err = msg.Err
		return m, nil
	}

	m.CurrentArticle = msg.Article
	m.ArticleContent = msg.Article.Content

	// Mark article as read
	if m.State != nil {
		m.State.MarkAsRead(msg.Article.URL)
		// Update the read status in the feed list
		if m.CurrentFeed < len(m.Feeds) {
			for i := range m.Feeds[m.CurrentFeed].Item {
				if m.Feeds[m.CurrentFeed].Item[i].Link == msg.Article.URL {
					m.Feeds[m.CurrentFeed].Item[i].Read = true
					break
				}
			}
		}
	}

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

	// Save state after marking as read
	return m, SaveState(m.State)
}

func handleStateLoad(m *Model, msg StateLoadMsg) (*Model, tea.Cmd) {
	if msg.Err != nil {
		// If state loading fails, just use empty state
		m.State = &storage.AppState{
			ReadArticles: make(map[string]bool),
			FeedURLs:     []storage.FeedConfig{},
		}
		return m, nil
	}

	m.State = msg.State

	// Mark articles as read based on loaded state
	for i := range m.Feeds {
		for j := range m.Feeds[i].Item {
			m.Feeds[i].Item[j].Read = m.State.IsRead(m.Feeds[i].Item[j].Link)
		}
	}

	return m, nil
}

func handleStateSave(m *Model, msg StateSaveMsg) (*Model, tea.Cmd) {
	if msg.Err != nil {
		m.Err = msg.Err
	}
	return m, nil
}

func handleConfigLoad(m *Model, msg ConfigLoadMsg) (*Model, tea.Cmd) {
	if msg.Err != nil {
		m.Err = msg.Err
		return m, nil
	}

	m.Config = msg.Config

	// Load all feeds from config
	var cmds []tea.Cmd
	for _, feedConfig := range msg.Config.Feeds {
		cmds = append(cmds, LoadFeed(feedConfig.URL))
	}

	// Batch all feed load commands
	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func handleFeedsLoaded(m *Model, msg FeedsLoadedMsg) (*Model, tea.Cmd) {
	// All feeds have been requested to load
	// Individual FeedLoadMsg messages will arrive as they complete
	return m, nil
}

func handleFeedAdded(m *Model, msg FeedAddedMsg) (*Model, tea.Cmd) {
	if msg.Err != nil {
		m.Err = msg.Err
		return m, nil
	}

	// Feed added successfully
	m.Err = nil
	return m, nil
}

func handleFeedDeleted(m *Model, msg FeedDeletedMsg) (*Model, tea.Cmd) {
	if msg.Err != nil {
		m.Err = msg.Err
		return m, nil
	}

	// Feed deleted successfully
	// Also remove the corresponding loaded feed if it exists
	if msg.Index < len(m.Feeds) {
		m.Feeds = append(m.Feeds[:msg.Index], m.Feeds[msg.Index+1:]...)
	}

	// Adjust cursor if needed
	if m.Cursor >= len(m.Config.Feeds) && m.Cursor > 0 {
		m.Cursor--
	}

	m.Err = nil
	return m, nil
}

func handleFeedUpdated(m *Model, msg FeedUpdatedMsg) (*Model, tea.Cmd) {
	if msg.Err != nil {
		m.Err = msg.Err
		return m, nil
	}

	// Feed updated successfully
	// Reload the feed
	m.Err = nil
	return m, LoadFeed(msg.Feed.URL)
}

func handleConfigSaved(m *Model, msg ConfigSavedMsg) (*Model, tea.Cmd) {
	if msg.Err != nil {
		m.Err = msg.Err
		return m, nil
	}

	m.Err = nil
	return m, nil
}

func handleClipboardPaste(m *Model, msg ClipboardPasteMsg) (*Model, tea.Cmd) {
	if msg.Err != nil {
		m.Err = msg.Err
		return m, nil
	}

	// Paste into the current field
	if m.AddingFeed {
		switch m.AddFeedField {
		case "url":
			m.AddFeedURL += msg.Content
		case "category":
			m.AddFeedCat += msg.Content
		case "tags":
			m.AddFeedTags += msg.Content
		}
	} else if m.EditingFeed {
		m.EditValue += msg.Content
	}

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
