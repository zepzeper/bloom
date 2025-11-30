package tui

import (
	"bloom/internal/feed"
	"bloom/internal/storage"
	"bloom/internal/tui/utils"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
)

// handleKeyMsg handles keyboard input for the model
func handleKeyMsg(m *Model, msg tea.KeyMsg) (*Model, tea.Cmd) {
	if m.Loading {
		return m, nil
	}

	// Handle feed management input
	if m.CurrentView == "manage" {
		return handleFeedManagementKeys(m, msg)
	}

	// Handle scrolling and cursor movement in content view first
	if m.CurrentView == "content" {
		switch msg.String() {
		case "j", "down":
			// Move cursor down or scroll
			return handleContentDown(m)
		case "k", "up":
			// Move cursor up or scroll
			return handleContentUp(m)
		case "h", "left":
			// Move cursor left
			return handleContentLeft(m)
		case "l", "right":
			// Move cursor right
			return handleContentRight(m)
		case "ctrl+u":
			return scrollPageUp(m)
		case "ctrl+d":
			return scrollPageDown(m)
		case "g":
			// Go to top (like vim gg)
			m.ScrollOffset = 0
			m.CursorY = 0
			m.CursorX = 0
			return m, nil
		case "G":
			// Go to bottom (like vim G)
			if len(m.ArticleLines) > 0 {
				m.ScrollOffset = len(m.ArticleLines) - 1
				m.CursorY = 0
				m.CursorX = 0
			}
			return m, nil
		case "o":
			// Open link under cursor
			return openLinkUnderCursor(m)
		case "c":
			// Copy link under cursor
			return copyLinkUnderCursor(m)
		}
	}

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "j", "down":
		return handleDown(m)
	case "k", "up":
		return handleUp(m)
	case "enter":
		return handleEnter(m)
	case "esc":
		return handleEscape(m)
	case "m":
		// Open feed manager (from landing or feed view)
		if m.CurrentView == "landing" || m.CurrentView == "feed" {
			m.CurrentView = "manage"
			m.Cursor = 0
			return m, nil
		}
		// Toggle read/unread status (in articles view)
		if m.CurrentView == "articles" {
			return toggleReadStatus(m)
		}
		return m, nil
	case "s":
		// Save state manually
		if m.State != nil {
			return m, SaveState(m.State)
		}
		return m, nil
	case "f":
		// Go to feeds view (from landing or manage)
		if m.CurrentView == "landing" {
			m.CurrentView = "feed"
			m.Cursor = 0
			return m, nil
		}
		// Open feed manager (from feed view)
		if m.CurrentView == "feed" {
			m.CurrentView = "manage"
			m.Cursor = 0
			return m, nil
		}
	}

	return m, nil
}

func handleDown(m *Model) (*Model, tea.Cmd) {
	switch m.CurrentView {
	case "feed":
		feedCount := len(m.Config.Feeds)
		if m.Config == nil {
			feedCount = 0
		}
		if m.CurrentFeed < feedCount-1 {
			m.CurrentFeed++
		}
	case "articles":
		if m.Config != nil && m.CurrentFeed < len(m.Config.Feeds) {
			// Find loaded feed
			for i := range m.Feeds {
				if m.Feeds[i].FeedURL == m.Config.Feeds[m.CurrentFeed].URL {
					if m.Cursor < len(m.Feeds[i].Item)-1 {
						m.Cursor++
					}
					break
				}
			}
		}
	}
	return m, nil
}

func handleUp(m *Model) (*Model, tea.Cmd) {
	switch m.CurrentView {
	case "feed":
		if m.CurrentFeed > 0 {
			m.CurrentFeed--
		}
	case "articles":
		if m.Cursor > 0 {
			m.Cursor--
		}
	}
	return m, nil
}

func handleEnter(m *Model) (*Model, tea.Cmd) {
	switch m.CurrentView {
	case "feed":
		if m.Config != nil && m.CurrentFeed < len(m.Config.Feeds) {
			// Check if feed is loaded and has articles
			for i := range m.Feeds {
				if m.Feeds[i].FeedURL == m.Config.Feeds[m.CurrentFeed].URL {
					if len(m.Feeds[i].Item) > 0 {
						m.CurrentView = "articles"
						m.Cursor = 0
					}
					break
				}
			}
		}
		return m, nil

	case "articles":
		if m.Config != nil && m.CurrentFeed < len(m.Config.Feeds) {
			// Find loaded feed
			for i := range m.Feeds {
				if m.Feeds[i].FeedURL == m.Config.Feeds[m.CurrentFeed].URL {
					feed := m.Feeds[i]
					if m.Cursor < len(feed.Item) {
						item := feed.Item[m.Cursor]
						m.Loading = true
						return m, LoadArticle(m.Fetcher, item.Link)
					}
					break
				}
			}
		}
		return m, nil

	case "content":
		return m, nil
	}

	return m, nil
}

func handleEscape(m *Model) (*Model, tea.Cmd) {
	switch m.CurrentView {
	case "feed":
		m.CurrentView = "landing"
		m.Cursor = 0
	case "articles":
		m.CurrentView = "feed"
		m.Cursor = 0
	case "content":
		m.CurrentView = "articles"
		m.ArticleContent = ""
		m.CurrentArticle = feed.Article{}
		m.CursorX = 0
		m.CursorY = 0
	case "manage":
		m.CurrentView = "landing"
		m.Cursor = 0
	}
	return m, nil
}

// toggleReadStatus toggles the read/unread status of the current article
func toggleReadStatus(m *Model) (*Model, tea.Cmd) {
	if m.CurrentView != "articles" {
		return m, nil
	}

	if m.State == nil {
		return m, nil
	}

	if m.Config == nil || m.CurrentFeed >= len(m.Config.Feeds) {
		return m, nil
	}

	// Find loaded feed by FeedURL
	feedConfig := m.Config.Feeds[m.CurrentFeed]
	var loadedFeed *feed.Channel
	for i := range m.Feeds {
		if m.Feeds[i].FeedURL == feedConfig.URL {
			loadedFeed = &m.Feeds[i]
			break
		}
	}

	if loadedFeed == nil || m.Cursor >= len(loadedFeed.Item) {
		return m, nil
	}

	item := &loadedFeed.Item[m.Cursor]
	item.Read = !item.Read

	if item.Read {
		m.State.MarkAsRead(item.Link)
	} else {
		// Remove from read articles
		delete(m.State.ReadArticles, item.Link)
	}

	return m, SaveState(m.State)
}

// handleFeedManagementKeys handles keyboard input in feed management view
func handleFeedManagementKeys(m *Model, msg tea.KeyMsg) (*Model, tea.Cmd) {
	// Handle adding a new feed
	if m.AddingFeed {
		return handleAddFeedKeys(m, msg)
	}

	// Handle editing an existing feed
	if m.EditingFeed {
		return handleEditFeedKeys(m, msg)
	}

	// Normal feed management navigation
	switch msg.String() {
	case "esc":
		m.CurrentView = "feed"
		m.Cursor = 0
		return m, nil
	case "j", "down":
		if m.Cursor < len(m.Config.Feeds)-1 {
			m.Cursor++
		}
		return m, nil
	case "k", "up":
		if m.Cursor > 0 {
			m.Cursor--
		}
		return m, nil
	case "a":
		// Start adding a new feed
		m.AddingFeed = true
		m.AddFeedURL = ""
		m.AddFeedCat = ""
		m.AddFeedTags = ""
		m.AddFeedField = "url"
		return m, nil
	case "e":
		// Start editing current feed
		if m.Cursor < len(m.Config.Feeds) {
			feed := m.Config.Feeds[m.Cursor]
			m.EditingFeed = true
			m.EditField = "url"
			m.EditValue = feed.URL
			return m, nil
		}
		return m, nil
	case "d":
		// Delete current feed
		if m.Cursor < len(m.Config.Feeds) {
			return m, DeleteFeedFromConfig(m.Config, m.Cursor)
		}
		return m, nil
	case "r":
		// Reload feeds from config
		return m, LoadConfig()
	}

	return m, nil
}

// handleAddFeedKeys handles keyboard input when adding a new feed
func handleAddFeedKeys(m *Model, msg tea.KeyMsg) (*Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Cancel adding
		m.AddingFeed = false
		m.AddFeedURL = ""
		m.AddFeedCat = ""
		m.AddFeedTags = ""
		return m, nil
	case "tab":
		// Move to next field
		switch m.AddFeedField {
		case "url":
			m.AddFeedField = "category"
		case "category":
			m.AddFeedField = "tags"
		case "tags":
			m.AddFeedField = "url"
		}
		return m, nil
	case "enter":
		// Save new feed
		if m.AddFeedURL == "" {
			m.Err = fmt.Errorf("URL is required")
			return m, nil
		}

		tags := []string{}
		if m.AddFeedTags != "" {
			for _, tag := range strings.Split(m.AddFeedTags, ",") {
				tags = append(tags, strings.TrimSpace(tag))
			}
		}

		newFeed := storage.FeedConfig{
			URL:      m.AddFeedURL,
			Category: m.AddFeedCat,
			Tags:     tags,
		}

		m.AddingFeed = false
		m.AddFeedURL = ""
		m.AddFeedCat = ""
		m.AddFeedTags = ""
		
		return m, tea.Batch(
			AddFeedToConfig(m.Config, newFeed),
			LoadFeed(newFeed.URL),
		)
	case "ctrl+v":
		// Paste from clipboard
		return m, PasteFromClipboard()
	case "backspace":
		// Delete character from current field
		switch m.AddFeedField {
		case "url":
			if len(m.AddFeedURL) > 0 {
				m.AddFeedURL = m.AddFeedURL[:len(m.AddFeedURL)-1]
			}
		case "category":
			if len(m.AddFeedCat) > 0 {
				m.AddFeedCat = m.AddFeedCat[:len(m.AddFeedCat)-1]
			}
		case "tags":
			if len(m.AddFeedTags) > 0 {
				m.AddFeedTags = m.AddFeedTags[:len(m.AddFeedTags)-1]
			}
		}
		return m, nil
	default:
		// Add character to current field
		if len(msg.String()) == 1 {
			switch m.AddFeedField {
			case "url":
				m.AddFeedURL += msg.String()
			case "category":
				m.AddFeedCat += msg.String()
			case "tags":
				m.AddFeedTags += msg.String()
			}
		}
		return m, nil
	}
}

// handleEditFeedKeys handles keyboard input when editing a feed
func handleEditFeedKeys(m *Model, msg tea.KeyMsg) (*Model, tea.Cmd) {
	if m.Cursor >= len(m.Config.Feeds) {
		m.EditingFeed = false
		return m, nil
	}

	feed := m.Config.Feeds[m.Cursor]

	switch msg.String() {
	case "esc":
		// Cancel editing
		m.EditingFeed = false
		m.EditField = ""
		m.EditValue = ""
		return m, nil
	case "tab":
		// Move to next field and save current
		switch m.EditField {
		case "url":
			feed.URL = m.EditValue
			m.EditField = "category"
			m.EditValue = feed.Category
		case "category":
			feed.Category = m.EditValue
			m.EditField = "tags"
			m.EditValue = strings.Join(feed.Tags, ", ")
		case "tags":
			// Parse tags
			tags := []string{}
			if m.EditValue != "" {
				for _, tag := range strings.Split(m.EditValue, ",") {
					tags = append(tags, strings.TrimSpace(tag))
				}
			}
			feed.Tags = tags
			m.EditField = "url"
			m.EditValue = feed.URL
		}
		m.Config.Feeds[m.Cursor] = feed
		return m, nil
	case "enter":
		// Save changes
		switch m.EditField {
		case "url":
			feed.URL = m.EditValue
		case "category":
			feed.Category = m.EditValue
		case "tags":
			tags := []string{}
			if m.EditValue != "" {
				for _, tag := range strings.Split(m.EditValue, ",") {
					tags = append(tags, strings.TrimSpace(tag))
				}
			}
			feed.Tags = tags
		}

		m.EditingFeed = false
		m.EditField = ""
		m.EditValue = ""
		
		return m, UpdateFeedInConfig(m.Config, m.Cursor, feed)
	case "ctrl+v":
		// Paste from clipboard
		return m, PasteFromClipboard()
	case "backspace":
		// Delete character
		if len(m.EditValue) > 0 {
			m.EditValue = m.EditValue[:len(m.EditValue)-1]
		}
		return m, nil
	default:
		// Add character
		if len(msg.String()) == 1 {
			m.EditValue += msg.String()
		}
		return m, nil
	}
}

// Content view cursor movement handlers
func handleContentDown(m *Model) (*Model, tea.Cmd) {
	if m.CurrentView != "content" {
		return m, nil
	}

	// Full screen: height - header (2) - status (1)
	visibleHeight := m.Height - 3
	if visibleHeight < 1 {
		visibleHeight = 10
	}

	// Move cursor down within visible area
	if m.CursorY < visibleHeight-1 {
		m.CursorY++
		// Adjust X to not exceed line length
		actualLine := m.ScrollOffset + m.CursorY
		if actualLine >= 0 && actualLine < len(m.ArticleLines) {
			// Get display width of line (accounting for ANSI codes)
			line := m.ArticleLines[actualLine]
			lineLen := runewidth.StringWidth(utils.StripANSI(line))
			if m.CursorX > lineLen {
				m.CursorX = lineLen
			}
		}
	} else {
		// Scroll down if at bottom of visible area
		return scrollDown(m)
	}
	return m, nil
}

func handleContentUp(m *Model) (*Model, tea.Cmd) {
	if m.CurrentView != "content" {
		return m, nil
	}

	// Move cursor up
	if m.CursorY > 0 {
		m.CursorY--
		// Adjust X to not exceed line length
		actualLine := m.ScrollOffset + m.CursorY
		if actualLine >= 0 && actualLine < len(m.ArticleLines) {
			line := m.ArticleLines[actualLine]
			lineLen := runewidth.StringWidth(utils.StripANSI(line))
			if m.CursorX > lineLen {
				m.CursorX = lineLen
			}
		}
	} else if m.ScrollOffset > 0 {
		// Scroll up if at top of visible area
		return scrollUp(m)
	}
	return m, nil
}

func handleContentLeft(m *Model) (*Model, tea.Cmd) {
	if m.CurrentView != "content" {
		return m, nil
	}

	if m.CursorX > 0 {
		m.CursorX--
	}
	return m, nil
}

func handleContentRight(m *Model) (*Model, tea.Cmd) {
	if m.CurrentView != "content" {
		return m, nil
	}

	actualLine := m.ScrollOffset + m.CursorY
	if actualLine >= 0 && actualLine < len(m.ArticleLines) {
		line := m.ArticleLines[actualLine]
		lineLen := runewidth.StringWidth(utils.StripANSI(line))
		if m.CursorX < lineLen {
			m.CursorX++
		}
	}
	return m, nil
}

func openLinkUnderCursor(m *Model) (*Model, tea.Cmd) {
	link := findLinkAtPosition(m)
	if link == nil {
		return m, nil
	}

	// Return a command to open the link
	return m, OpenLink(link.URL)
}

func copyLinkUnderCursor(m *Model) (*Model, tea.Cmd) {
	link := findLinkAtPosition(m)
	if link == nil {
		return m, nil
	}

	// Return a command to copy the link
	return m, CopyLink(link.URL)
}
