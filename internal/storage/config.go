package storage

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// FeedConfig represents a feed configuration
type FeedConfig struct {
	URL      string
	Category string
	Tags     []string
}

// Config represents the application configuration
type Config struct {
	Feeds              []FeedConfig `json:"feeds"`
	AutoSave           bool         `json:"auto_save"`
	MarkReadOnView     bool         `json:"mark_read_on_view"`
	DefaultCategory    string       `json:"default_category"`
	RefreshIntervalMin int          `json:"refresh_interval_min"`
}

// LoadConfig loads the configuration from ~/.config/bloom/config.json
func LoadConfig() (*Config, error) {
	// Get config directory path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return DefaultConfig(), nil // Return default config if can't find home
	}

	configDir := filepath.Join(homeDir, ".config", "bloom")
	configPath := filepath.Join(configDir, "config.json")

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config file
		defaultConfig := DefaultConfig()
		err := SaveConfig(defaultConfig)
		if err != nil {
			return defaultConfig, nil // Return default even if save fails
		}
		return defaultConfig, nil
	}

	// Read file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// Ensure feeds slice is initialized
	if config.Feeds == nil {
		config.Feeds = []FeedConfig{}
	}

	// Normalize feed URLs (add https:// if missing)
	needsSave := false
	for i := range config.Feeds {
		normalized := normalizeFeedURL(config.Feeds[i].URL)
		if normalized != config.Feeds[i].URL {
			config.Feeds[i].URL = normalized
			needsSave = true
		}
	}
	
	// Save config with normalized URLs if any changed
	if needsSave {
		if err := SaveConfig(&config); err != nil {
			// Log error but don't fail - config is still valid
			fmt.Printf("Warning: failed to save normalized URLs: %v\n", err)
		}
	}

	return &config, nil
}

// normalizeFeedURL normalizes a feed URL by adding https:// if no protocol is present
func normalizeFeedURL(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return rawURL
	}
	
	// Check if URL already has a protocol
	parsed, err := url.Parse(rawURL)
	if err != nil {
		// If parsing fails, try adding https://
		return "https://" + rawURL
	}
	
	// If no scheme, add https://
	if parsed.Scheme == "" {
		return "https://" + rawURL
	}
	
	return rawURL
}

// SaveConfig saves the configuration to ~/.config/bloom/config.json
func SaveConfig(config *Config) error {
	// Get config directory path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ".config", "bloom")
	configPath := filepath.Join(configDir, "config.json")

	// Create config directory if it doesn't exist
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	// Write to file
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Feeds: []FeedConfig{
			{
				URL:      "https://mitchellh.com/feed.xml",
				Category: "Tech",
				Tags:     []string{"golang", "infrastructure"},
			},
		},
		AutoSave:           true,
		MarkReadOnView:     true,
		DefaultCategory:    "",
		RefreshIntervalMin: 60,
	}
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "bloom", "config.json"), nil
}

