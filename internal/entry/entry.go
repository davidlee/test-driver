// Package entry formats text bodies for log entries.
package entry

import "strings"

// Format prepares body text for appending to a log file.
// If task is true, each non-empty line is prefixed with a Markdown checkbox.
// Returns empty string for empty/whitespace-only input.
func Format(body string, task bool) string {
	trimmed := strings.TrimSpace(body)
	if trimmed == "" {
		return ""
	}

	if !task {
		return trimmed + "\n"
	}

	lines := strings.Split(trimmed, "\n")
	var out []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		out = append(out, "- [ ] "+line)
	}
	return strings.Join(out, "\n") + "\n"
}
