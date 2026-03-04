package logfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"im/internal/config"
	"im/internal/entry"
)

// Appender appends formatted entries to daily log files.
type Appender struct {
	clock    func() time.Time
	strategy config.TimestampRounding
}

// NewAppender creates an Appender. The clock function is called to get the
// current time (inject time.Now in production, a fixed clock in tests).
func NewAppender(clock func() time.Time, cfg config.Config) *Appender {
	return &Appender{
		clock:    clock,
		strategy: cfg.TimestampRounding,
	}
}

// Append writes body to the daily log file in dir.
// If task is true, each line is prefixed with a Markdown checkbox.
func (a *Appender) Append(dir, body string, task bool) error {
	formatted := entry.Format(body, task)
	if formatted == "" {
		return nil
	}

	now := a.clock()
	filename := now.Format("2006-01-02") + ".md"
	path := filepath.Join(dir, filename)

	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("creating log directory: %w", err)
	}

	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading log file: %w", err)
	}

	var buf strings.Builder
	isNew := len(existing) == 0

	if isNew {
		buf.WriteString(dateHeading(now))
		buf.WriteString("\n")
	}

	parsedLast, hasLast := ParseLastTimestamp(string(existing))
	// Reconstruct parsed time on today's date for correct gap calculation.
	last := time.Date(now.Year(), now.Month(), now.Day(),
		parsedLast.Hour(), parsedLast.Minute(), 0, 0, now.Location())
	emit, display := ShouldEmitHeading(a.strategy, now, last, hasLast)

	if emit {
		buf.WriteString("\n")
		fmt.Fprintf(&buf, "%s %s\n", SubheadingLevel, display.Format("15:04"))
	}

	buf.WriteString("\n")
	buf.WriteString(formatted)

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("opening log file: %w", err)
	}

	_, writeErr := f.WriteString(buf.String())
	closeErr := f.Close()
	if writeErr != nil {
		return fmt.Errorf("writing to log file: %w", writeErr)
	}
	if closeErr != nil {
		return fmt.Errorf("closing log file: %w", closeErr)
	}

	return nil
}

// dateHeading returns the level-1 heading for a daily log file.
func dateHeading(t time.Time) string {
	return fmt.Sprintf("# %s", t.Format("Monday, 2 January 2006"))
}
