package tui

import tea "github.com/charmbracelet/bubbletea"

// Init initializes the model (bubbletea interface)
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		LoadState(),
		LoadConfig(),
	)
}
