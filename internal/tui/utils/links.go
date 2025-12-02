package utils

import (
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"
)

// Link represents a link found in the article
type Link struct {
	Text string
	URL  string
	Line int // Line number in articleLines
	Start int // Start position in line
	End   int // End position in line
}

// ParseLinksFromMarkdown extracts links from markdown content
func ParseLinksFromMarkdown(markdown string) []Link {
	var links []Link
	lines := strings.Split(markdown, "\n")

	// Pattern for markdown links: [text](url)
	linkPattern := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

	for lineNum, line := range lines {
		matches := linkPattern.FindAllStringSubmatchIndex(line, -1)
		for _, match := range matches {
			if len(match) >= 6 {
				textStart := match[2]
				textEnd := match[3]
				urlStart := match[4]
				urlEnd := match[5]

				text := line[textStart:textEnd]
				url := line[urlStart:urlEnd]

				links = append(links, Link{
					Text:  text,
					URL:   url,
					Line:  lineNum,
					Start: match[0], // Start of entire match
					End:   match[1], // End of entire match
				})
			}
		}

		// Also check for plain URLs
		urlPattern := regexp.MustCompile(`https?://[^\s]+`)
		urlMatches := urlPattern.FindAllStringIndex(line, -1)
		for _, match := range urlMatches {
			url := line[match[0]:match[1]]
			links = append(links, Link{
				Text:  url,
				URL:   url,
				Line:  lineNum,
				Start: match[0],
				End:   match[1],
			})
		}
	}

	return links
}

// ParseLinksFromRenderedContent extracts links from rendered content (with ANSI codes)
func ParseLinksFromRenderedContent(rendered string) []Link {
	var links []Link
	lines := strings.Split(rendered, "\n")

	// Pattern for URLs (http/https) - more permissive to catch URLs with various characters
	urlPattern := regexp.MustCompile(`https?://[^\s\x1b\[\]()]+`)

	for lineNum, line := range lines {
		// Find all URL matches
		matches := urlPattern.FindAllStringIndex(line, -1)
		for _, match := range matches {
			// Extract the actual URL (may have ANSI codes)
			urlWithCodes := line[match[0]:match[1]]
			// Strip ANSI codes to get clean URL
			url := StripANSI(urlWithCodes)
			
			// Skip if URL is empty after stripping
			if url == "" {
				continue
			}

			// Calculate display position (accounting for ANSI codes)
			// We need to calculate the display width of everything before the match
			beforeMatch := line[:match[0]]
			displayStart := runewidth.StringWidth(StripANSI(beforeMatch))
			displayEnd := displayStart + runewidth.StringWidth(url)

			links = append(links, Link{
				Text:  url,
				URL:   url,
				Line:  lineNum,
				Start: displayStart,
				End:   displayEnd - 1, // Make End inclusive (cursor can be at End position)
			})
		}
	}

	return links
}

// StripANSI removes ANSI escape codes from a string
func StripANSI(s string) string {
	ansiPattern := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiPattern.ReplaceAllString(s, "")
}

// FindURLInRenderedLine finds the URL text in a rendered line (with ANSI codes)
func FindURLInRenderedLine(line, url string) string {
	// Try to find the URL in the line - it might have ANSI codes around it
	// Look for the URL pattern
	urlPattern := regexp.MustCompile(regexp.QuoteMeta(url))
	match := urlPattern.FindString(line)
	if match != "" {
		return match
	}

	// If not found, try finding it with potential ANSI codes
	// Look for URL with optional ANSI codes before/after
	ansiURLPattern := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]*` + regexp.QuoteMeta(url) + `\x1b\[[0-9;]*[a-zA-Z]*`)
	match = ansiURLPattern.FindString(line)
	if match != "" {
		return match
	}

	return ""
}
