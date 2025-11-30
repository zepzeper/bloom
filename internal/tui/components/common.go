package components

import (
	"bloom/internal/tui/styles"
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func RenderErrorView(err error, width int) string {
	// Simple error display - no fancy boxes
	content := styles.ErrorStyle().Render(fmt.Sprintf("Error: %v", err))
	status := styles.RenderStatusBar("ERROR", "", "Press q to quit", width)
	return lipgloss.JoinVertical(lipgloss.Left, content, status)
}

func RenderLoadingView(width int) string {
	// Simple loading text
	content := styles.LoadingStyle().Render("Loading...")
	status := styles.RenderStatusBar("Loading", "", "Press q to quit", width)
	return lipgloss.JoinVertical(lipgloss.Left, content, status)
}
