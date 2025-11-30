package handlers

import (
	"bloom/internal/feed"
	"bloom/internal/tui"
	"bloom/internal/tui/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
)

// HandleKeyMsg handles keyboard input for the model
func HandleKeyMsg(m *tui.Model, msg tea.KeyMsg) (*tui.Model, tea.Cmd) {
	if m.Loading {
		return m, nil
	}

	// Handle scrolling and cursor movement in content view first
	if m.CurrentView == "content" {
		switch msg.String() {
		case "j", "down":
			// Move cursor down or scroll
			return HandleContentDown(m)
		case "k", "up":
			// Move cursor up or scroll
			return HandleContentUp(m)
		case "h", "left":
			// Move cursor left
			return HandleContentLeft(m)
		case "l", "right":
			// Move cursor right
			return HandleContentRight(m)
		case "ctrl+u":
			return ScrollPageUp(m)
		case "ctrl+d":
			return ScrollPageDown(m)
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
			return OpenLinkUnderCursor(m)
		case "c":
			// Copy link under cursor
			return CopyLinkUnderCursor(m)
		}
	}

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "j", "down":
		return HandleDown(m)
	case "k", "up":
		return HandleUp(m)
	case "enter":
		return HandleEnter(m)
	case "esc":
		return HandleEscape(m)
	}

	return m, nil
}

func HandleDown(m *tui.Model) (*tui.Model, tea.Cmd) {
	switch m.CurrentView {
	case "feed":
		if m.CurrentFeed < len(m.Feeds)-1 {
			m.CurrentFeed++
		}
	case "articles":
		if len(m.Feeds) > 0 && m.CurrentFeed < len(m.Feeds) {
			if m.Cursor < len(m.Feeds[m.CurrentFeed].Item)-1 {
				m.Cursor++
			}
		}
	}
	return m, nil
}

func HandleUp(m *tui.Model) (*tui.Model, tea.Cmd) {
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

func HandleEnter(m *tui.Model) (*tui.Model, tea.Cmd) {
	switch m.CurrentView {
	case "feed":
		if len(m.Feeds) > 0 && m.CurrentFeed < len(m.Feeds) {
			if len(m.Feeds[m.CurrentFeed].Item) > 0 {
				m.CurrentView = "articles"
				m.Cursor = 0
			}
		}
		return m, nil

	case "articles":
		if len(m.Feeds) > 0 && m.CurrentFeed < len(m.Feeds) {
			feed := m.Feeds[m.CurrentFeed]
			if m.Cursor < len(feed.Item) {
				item := feed.Item[m.Cursor]
				m.Loading = true
				return m, tui.LoadArticle(m.Fetcher, item.Link)
			}
		}
		return m, nil

	case "content":
		return m, nil
	}

	return m, nil
}

func HandleEscape(m *tui.Model) (*tui.Model, tea.Cmd) {
	switch m.CurrentView {
	case "articles":
		m.CurrentView = "feed"
		m.Cursor = 0
	case "content":
		m.CurrentView = "articles"
		m.ArticleContent = ""
		m.CurrentArticle = feed.Article{}
		m.CursorX = 0
		m.CursorY = 0
	}
	return m, nil
}

// Content view cursor movement handlers
func HandleContentDown(m *tui.Model) (*tui.Model, tea.Cmd) {
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
		return ScrollDown(m)
	}
	return m, nil
}

func HandleContentUp(m *tui.Model) (*tui.Model, tea.Cmd) {
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
		return ScrollUp(m)
	}
	return m, nil
}

func HandleContentLeft(m *tui.Model) (*tui.Model, tea.Cmd) {
	if m.CurrentView != "content" {
		return m, nil
	}

	if m.CursorX > 0 {
		m.CursorX--
	}
	return m, nil
}

func HandleContentRight(m *tui.Model) (*tui.Model, tea.Cmd) {
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

func OpenLinkUnderCursor(m *tui.Model) (*tui.Model, tea.Cmd) {
	link := FindLinkAtPosition(m)
	if link == nil {
		return m, nil
	}

	// Return a command to open the link
	return m, tui.OpenLink(link.URL)
}

func CopyLinkUnderCursor(m *tui.Model) (*tui.Model, tea.Cmd) {
	link := FindLinkAtPosition(m)
	if link == nil {
		return m, nil
	}

	// Return a command to copy the link
	return m, tui.CopyLink(link.URL)
}
