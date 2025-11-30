package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type AppState struct {
	ReadArticles map[string]bool
	FeedURLs     []FeedConfig
	LastSync     time.Time
}

type FeedConfig struct {
	URL      string
	Category string
	Tags     []string
}

func (s *AppState) MarkAsRead(articleURL string) {
	s.ReadArticles[articleURL] = true
}

func (s *AppState) IsRead(articleURL string) bool {
	return s.ReadArticles[articleURL]
}

func LoadState() (*AppState, error) {
	// Get config directory path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return NewAppState(), nil // Return empty state if can't find home
	}

	configDir := filepath.Join(homeDir, ".config", "bloom")
	statePath := filepath.Join(configDir, "state.json")

	// Check if file exists
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		// Return new state if file doesn't exist
		return NewAppState(), nil
	}

	// Read file
	data, err := os.ReadFile(statePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %v", err)
	}

	var state AppState
	err = json.Unmarshal(data, &state)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal state: %v", err)
	}

	// Initialize maps if nil
	if state.ReadArticles == nil {
		state.ReadArticles = make(map[string]bool)
	}
	if state.FeedURLs == nil {
		state.FeedURLs = []FeedConfig{}
	}

	return &state, nil
}

func SaveState(state *AppState) error {
	// Get config directory path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ".config", "bloom")
	statePath := filepath.Join(configDir, "state.json")

	// Create config directory if it doesn't exist
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// Marshal state to JSON
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %v", err)
	}

	// Write to file
	err = os.WriteFile(statePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write state file: %v", err)
	}

	return nil
}

// NewAppState creates a new empty AppState
func NewAppState() *AppState {
	return &AppState{
		ReadArticles: make(map[string]bool),
		FeedURLs:     []FeedConfig{},
		LastSync:     time.Now(),
	}
}
