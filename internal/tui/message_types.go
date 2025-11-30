package tui

import (
	"bloom/internal/feed"
	"bloom/internal/storage"
)

type StateLoadMsg struct {
	State *storage.AppState
	Err   error
}

type StateSaveMsg struct {
	Err error
}

// FeedLoadMsg is sent when a feed has been loaded
type FeedLoadMsg struct {
	Channel *feed.Channel
	Err     error
}

// ArticleLoadMsg is sent when an article has been loaded
type ArticleLoadMsg struct {
	Article feed.Article
	Err     error
}

// LinkOpenedMsg is sent when a link has been opened
type LinkOpenedMsg struct {
	URL string
	Err error
}

// LinkCopiedMsg is sent when a link has been copied
type LinkCopiedMsg struct {
	URL string
	Err error
}

// ConfigLoadMsg is sent when the config has been loaded
type ConfigLoadMsg struct {
	Config *storage.Config
	Err    error
}

// FeedsLoadedMsg is sent when all feeds from config have been loaded
type FeedsLoadedMsg struct {
	Count int
}

// FeedAddedMsg is sent when a feed has been added to config
type FeedAddedMsg struct {
	Feed storage.FeedConfig
	Err  error
}

// FeedDeletedMsg is sent when a feed has been deleted from config
type FeedDeletedMsg struct {
	Index int
	Err   error
}

// FeedUpdatedMsg is sent when a feed has been updated in config
type FeedUpdatedMsg struct {
	Index int
	Feed  storage.FeedConfig
	Err   error
}

// ConfigSavedMsg is sent when config has been saved
type ConfigSavedMsg struct {
	Err error
}

// ClipboardPasteMsg is sent when clipboard content has been read
type ClipboardPasteMsg struct {
	Content string
	Err     error
}
