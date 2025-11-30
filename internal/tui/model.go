package tui

import (
	"bloom/internal/feed"
	"bloom/internal/storage"
	"bloom/internal/tui/utils"
)

// Model represents the application state for the TUI
type Model struct {
	// State and Config
	State  *storage.AppState
	Config *storage.Config

	// View state
	CurrentView string
	Cursor      int

	// Feed data
	Feeds       []feed.Channel
	CurrentFeed int

	// Article data
	ArticleContent string
	CurrentArticle feed.Article
	ArticleLines   []string
	ArticleLinks   []utils.Link

	Categories      map[string]int
	CurrentCategory int
	ShowCategories  bool

	// Services
	Reader  *feed.Reader
	Fetcher *feed.ArticleFetcher

	// UI state
	Loading bool
	Err     error

	// Window dimensions
	Width  int
	Height int

	// Scroll state
	ScrollOffset int
	CursorX      int
	CursorY      int

	// Feed management state
	EditingFeed   bool
	EditField     string // "url", "category", "tags"
	EditValue     string
	AddingFeed    bool
	AddFeedURL    string
	AddFeedCat    string
	AddFeedTags   string
	AddFeedField  string // Current field being edited when adding
}

// NewModel creates and initializes a new Model
func NewModel() Model {
	return Model{
		State:           storage.NewAppState(),
		Config:          storage.DefaultConfig(),
		CurrentView:     "landing",
		Cursor:          0,
		Feeds:           []feed.Channel{},
		CurrentFeed:     0,
		ArticleContent:  "",
		ArticleLines:    []string{},
		ArticleLinks:    []utils.Link{},
		Categories:      map[string]int{},
		CurrentCategory: 0,
		ShowCategories:  false,
		Reader:          feed.NewReader(),
		Fetcher:         feed.NewArticleFetcher(),
		Loading:         false,
		Err:             nil,
		Width:           80,
		Height:          24,
		ScrollOffset:    0,
		CursorX:         0,
		CursorY:         0,
		EditingFeed:     false,
		EditField:       "",
		EditValue:       "",
		AddingFeed:      false,
		AddFeedURL:      "",
		AddFeedCat:      "",
		AddFeedTags:     "",
		AddFeedField:    "url",
	}
}
