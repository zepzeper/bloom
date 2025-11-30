package tui

import (
	"bloom/internal/feed"
	"bloom/internal/storage"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func LoadState() tea.Cmd {
	return func() tea.Msg {
		state, err := storage.LoadState()
		return StateLoadMsg{State: state, Err: err}
	}
}

func SaveState(state *storage.AppState) tea.Cmd {
	return func() tea.Msg {
		err := storage.SaveState(state)
		return StateSaveMsg{Err: err}
	}
}

// LoadConfig loads the application configuration
func LoadConfig() tea.Cmd {
	return func() tea.Msg {
		config, err := storage.LoadConfig()
		return ConfigLoadMsg{Config: config, Err: err}
	}
}

// LoadFeedsFromConfig loads all feeds from the config
func LoadFeedsFromConfig(config *storage.Config, reader *feed.Reader) tea.Cmd {
	return func() tea.Msg {
		// Load all feeds in sequence
		// In a real app, you might want to do this concurrently
		var cmds []tea.Cmd
		for _, feedConfig := range config.Feeds {
			cmds = append(cmds, LoadFeed(feedConfig.URL))
		}
		
		// For now, just return a message indicating we're done
		// The feeds will be loaded via individual FeedLoadMsg messages
		return FeedsLoadedMsg{Count: len(config.Feeds)}
	}
}

// LoadFeed loads an RSS feed from a URL
func LoadFeed(url string) tea.Cmd {
	return func() tea.Msg {
		reader := feed.NewReader()
		channel, err := reader.Read(url)
		return FeedLoadMsg{Channel: channel, Err: err}
	}
}

// LoadArticle loads an article from a URL
func LoadArticle(fetcher *feed.ArticleFetcher, url string) tea.Cmd {
	return func() tea.Msg {
		article, err := fetcher.Extract(url)
		return ArticleLoadMsg{Article: article, Err: err}
	}
}

// OpenLink opens a URL in the default browser
func OpenLink(url string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "linux":
			cmd = exec.Command("xdg-open", url)
		case "darwin":
			cmd = exec.Command("open", url)
		case "windows":
			cmd = exec.Command("cmd", "/c", "start", url)
		default:
			return LinkOpenedMsg{URL: url, Err: exec.ErrNotFound}
		}

		err := cmd.Run()
		return LinkOpenedMsg{URL: url, Err: err}
	}
}

// SaveConfig saves the configuration file
func SaveConfig(config *storage.Config) tea.Cmd {
	return func() tea.Msg {
		err := storage.SaveConfig(config)
		return ConfigSavedMsg{Err: err}
	}
}

// AddFeedToConfig adds a new feed to the configuration
func AddFeedToConfig(config *storage.Config, feedConfig storage.FeedConfig) tea.Cmd {
	return func() tea.Msg {
		config.Feeds = append(config.Feeds, feedConfig)
		err := storage.SaveConfig(config)
		if err != nil {
			return FeedAddedMsg{Feed: feedConfig, Err: err}
		}
		return FeedAddedMsg{Feed: feedConfig, Err: nil}
	}
}

// DeleteFeedFromConfig removes a feed from the configuration
func DeleteFeedFromConfig(config *storage.Config, index int) tea.Cmd {
	return func() tea.Msg {
		if index < 0 || index >= len(config.Feeds) {
			return FeedDeletedMsg{Index: index, Err: fmt.Errorf("invalid feed index")}
		}
		
		// Remove feed at index
		config.Feeds = append(config.Feeds[:index], config.Feeds[index+1:]...)
		err := storage.SaveConfig(config)
		return FeedDeletedMsg{Index: index, Err: err}
	}
}

// UpdateFeedInConfig updates a feed in the configuration
func UpdateFeedInConfig(config *storage.Config, index int, feedConfig storage.FeedConfig) tea.Cmd {
	return func() tea.Msg {
		if index < 0 || index >= len(config.Feeds) {
			return FeedUpdatedMsg{Index: index, Feed: feedConfig, Err: fmt.Errorf("invalid feed index")}
		}
		
		config.Feeds[index] = feedConfig
		err := storage.SaveConfig(config)
		return FeedUpdatedMsg{Index: index, Feed: feedConfig, Err: err}
	}
}

// PasteFromClipboard reads content from the system clipboard
func PasteFromClipboard() tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd
		var output []byte
		var err error

		switch runtime.GOOS {
		case "linux":
			// Try xclip first, then xsel
			cmd = exec.Command("xclip", "-selection", "clipboard", "-o")
			output, err = cmd.Output()
			if err != nil {
				cmd = exec.Command("xsel", "--clipboard", "--output")
				output, err = cmd.Output()
			}
		case "darwin":
			cmd = exec.Command("pbpaste")
			output, err = cmd.Output()
		case "windows":
			// Windows clipboard via PowerShell
			cmd = exec.Command("powershell", "-Command", "Get-Clipboard")
			output, err = cmd.Output()
		default:
			return ClipboardPasteMsg{Content: "", Err: exec.ErrNotFound}
		}

		if err != nil {
			return ClipboardPasteMsg{Content: "", Err: err}
		}

		return ClipboardPasteMsg{Content: strings.TrimSpace(string(output)), Err: nil}
	}
}

// CopyLink copies a URL to clipboard
func CopyLink(url string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd
		var err error

		switch runtime.GOOS {
		case "linux":
			// Try xclip first, then xsel
			cmd = exec.Command("xclip", "-selection", "clipboard")
			cmd.Stdin = strings.NewReader(url)
			err = cmd.Run()
			if err != nil {
				cmd = exec.Command("xsel", "--clipboard", "--input")
				cmd.Stdin = strings.NewReader(url)
				err = cmd.Run()
			}
		case "darwin":
			cmd = exec.Command("pbcopy")
			cmd.Stdin = strings.NewReader(url)
			err = cmd.Run()
		case "windows":
			// Windows clipboard via PowerShell
			cmd = exec.Command("powershell", "-Command", "Set-Clipboard", "-Value", url)
			err = cmd.Run()
			if err != nil {
				// Fallback to clip.exe with echo
				cmd = exec.Command("cmd", "/c", "echo", url+"|", "clip")
				err = cmd.Run()
			}
		default:
			return LinkCopiedMsg{URL: url, Err: exec.ErrNotFound}
		}

		return LinkCopiedMsg{URL: url, Err: err}
	}
}
