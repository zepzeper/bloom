package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update is the main update dispatcher (bubbletea interface)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var newModel *Model

	switch msg := msg.(type) {
	case tea.KeyMsg:
		newModel, cmd = handleKeyMsg(&m, msg)
		return *newModel, cmd

	case FeedLoadMsg:
		newModel, cmd = handleFeedLoad(&m, msg)
		return *newModel, cmd

	case ArticleLoadMsg:
		newModel, cmd = handleArticleLoad(&m, msg)
		return *newModel, cmd

	case LinkOpenedMsg:
		// Link opened - could show a message or do nothing
		if msg.Err != nil {
			m.Err = msg.Err
		}
		return m, nil

	case LinkCopiedMsg:
		// Link copied - could show a message or do nothing
		if msg.Err != nil {
			m.Err = msg.Err
		}
		return m, nil

	case StateLoadMsg:
		newModel, cmd = handleStateLoad(&m, msg)
		return *newModel, cmd

	case StateSaveMsg:
		newModel, cmd = handleStateSave(&m, msg)
		return *newModel, cmd

	case ConfigLoadMsg:
		newModel, cmd = handleConfigLoad(&m, msg)
		return *newModel, cmd

	case FeedsLoadedMsg:
		newModel, cmd = handleFeedsLoaded(&m, msg)
		return *newModel, cmd

	case FeedAddedMsg:
		newModel, cmd = handleFeedAdded(&m, msg)
		return *newModel, cmd

	case FeedDeletedMsg:
		newModel, cmd = handleFeedDeleted(&m, msg)
		return *newModel, cmd

	case FeedUpdatedMsg:
		newModel, cmd = handleFeedUpdated(&m, msg)
		return *newModel, cmd

	case ConfigSavedMsg:
		newModel, cmd = handleConfigSaved(&m, msg)
		return *newModel, cmd

	case ClipboardPasteMsg:
		newModel, cmd = handleClipboardPaste(&m, msg)
		return *newModel, cmd

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil
	}

	return m, nil
}
