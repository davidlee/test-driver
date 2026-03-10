// Package logfile manages daily log file creation, timestamp logic, and append.
package logfile

import (
	"regexp"
	"strconv"
	"time"
)

// SubheadingLevel is the Markdown heading level used for timestamp subheadings.
const SubheadingLevel = "##"

// timestampRe matches timestamp subheadings in both 24h ("## 14:30") and
// 12h ("## 2:30 PM") formats.
var timestampRe = regexp.MustCompile(
	`(?m)^` + SubheadingLevel + ` (\d{1,2}):(\d{2})(?:\s*(AM|PM))?\s*$`,
)

// ParseLastTimestamp extracts the time from the last timestamp heading in
// content. Handles both 24h and 12h formats. Returns the parsed time (on a
// zero date) and true, or zero time and false if no valid heading is found.
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

	h = convert12hTo24h(h, last[3])

	if h > 23 || m > 59 {
		return time.Time{}, false
	}

	return time.Date(0, 1, 1, h, m, 0, 0, time.UTC), true
}

// convert12hTo24h adjusts hour for AM/PM. If ampm is empty, returns h unchanged.
func convert12hTo24h(h int, ampm string) int {
	if ampm == "PM" && h != 12 {
		return h + 12
	}
	if ampm == "AM" && h == 12 {
		return 0
	}
	return h
}
