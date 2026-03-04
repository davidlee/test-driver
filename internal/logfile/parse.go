// Package logfile manages daily log file creation, timestamp logic, and append.
package logfile

import (
	"regexp"
	"strconv"
	"time"
)

// SubheadingLevel is the Markdown heading level used for timestamp subheadings.
const SubheadingLevel = "##"

// timestampRe matches a timestamp subheading like "## 14:30".
var timestampRe = regexp.MustCompile(`(?m)^` + SubheadingLevel + ` (\d{2}):(\d{2})\s*$`)

// ParseLastTimestamp extracts the time from the last ## HH:MM heading in content.
// Returns the parsed time (on a zero date) and true, or zero time and false if
// no valid heading is found.
func ParseLastTimestamp(content string) (time.Time, bool) {
	matches := timestampRe.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return time.Time{}, false
	}

	last := matches[len(matches)-1]
	h, err := strconv.Atoi(last[1])
	if err != nil {
		return time.Time{}, false
	}
	m, err := strconv.Atoi(last[2])
	if err != nil {
		return time.Time{}, false
	}

	if h > 23 || m > 59 {
		return time.Time{}, false
	}

	t := time.Date(0, 1, 1, h, m, 0, 0, time.UTC)
	return t, true
}
