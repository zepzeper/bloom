package utils

import (
	"html"
	"regexp"
)

// StripHTML removes HTML tags and decodes entities (used for feed descriptions)
func StripHTML(htmlContent string) string {
	// Remove HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	text := re.ReplaceAllString(htmlContent, "")

	// Decode HTML entities
	text = html.UnescapeString(text)

	// Clean up whitespace
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	return text
}

