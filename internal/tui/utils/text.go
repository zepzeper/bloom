package utils

import (
	"html"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
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

// InsertCursorAtPosition inserts a cursor character at the given display position
// while preserving all ANSI color codes in the line
func InsertCursorAtPosition(line string, cursorPos int) string {
	if cursorPos < 0 {
		cursorPos = 0
	}

	// Calculate display width to find where to insert cursor
	displayPos := 0
	var result strings.Builder
	inAnsi := false

	for i := 0; i < len(line); i++ {
		char := line[i]

		// Check for ANSI escape sequence start
		if char == '\x1b' && i+1 < len(line) && line[i+1] == '[' {
			inAnsi = true
			result.WriteByte(char)
			continue
		}

		if inAnsi {
			result.WriteByte(char)
			// Check for ANSI sequence end (letter a-z, A-Z)
			if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
				inAnsi = false
			}
			continue
		}

		// Regular character - check if we've reached cursor position
		r, size := utf8.DecodeRuneInString(line[i:])
		if r == utf8.RuneError {
			break
		}

		charWidth := runewidth.RuneWidth(r)
		if displayPos == cursorPos {
			// Insert cursor here - reverse the character
			cursorChar := "\x1b[7m" + string(r) + "\x1b[27m"
			result.WriteString(cursorChar)
			// Write the rest of the line
			if i+size < len(line) {
				result.WriteString(line[i+size:])
			}
			return result.String()
		}

		if displayPos+charWidth > cursorPos {
			// Cursor is in the middle of a wide character
			// Just insert reverse space
			cursorChar := "\x1b[7m \x1b[27m"
			result.WriteString(cursorChar)
			result.WriteRune(r)
			// Write the rest of the line
			if i+size < len(line) {
				result.WriteString(line[i+size:])
			}
			return result.String()
		}

		displayPos += charWidth
		result.WriteRune(r)
		i += size - 1
	}

	// Cursor is at or beyond end of line - add reverse space
	cursorChar := "\x1b[7m \x1b[27m"
	result.WriteString(cursorChar)
	return result.String()
}

