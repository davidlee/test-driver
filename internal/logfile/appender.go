package logfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"im/internal/config"
	"im/internal/entry"
	"im/internal/tmpl"
)

// Appender appends formatted entries to daily log files.
type Appender struct {
	clock        func() time.Time
	strategy     config.TimestampRounding
	timeFormat   config.TimeFormat
	templatePath string
	titleFormat  string
}

// NewAppender creates an Appender. The clock function is called to get the
// current time (inject time.Now in production, a fixed clock in tests).
func NewAppender(clock func() time.Time, cfg config.Config) *Appender {
	return &Appender{
		clock:        clock,
		strategy:     cfg.TimestampRounding,
		timeFormat:   cfg.TimeFormat,
		templatePath: config.DefaultTemplatePath(),
		titleFormat:  cfg.TitleFormat,
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

	content, err := a.buildContent(now, existing, formatted)
	if err != nil {
		return err
	}

	return appendFile(path, content)
}

func (a *Appender) buildContent(now time.Time, existing []byte, formatted string) (string, error) {
	var buf strings.Builder

	if len(existing) == 0 {
		header, err := a.renderNewFileHeader(now)
		if err != nil {
			return "", err
		}
		buf.WriteString(header)
	}

	parsedLast, hasLast := ParseLastTimestamp(string(existing))
	last := time.Date(now.Year(), now.Month(), now.Day(),
		parsedLast.Hour(), parsedLast.Minute(), 0, 0, now.Location())
	emit, display := ShouldEmitHeading(a.strategy, now, last, hasLast)

	if emit {
		buf.WriteString("\n")
		fmt.Fprintf(&buf, "%s %s\n", SubheadingLevel, FormatTime(display, a.timeFormat))
	}

	buf.WriteString("\n")
	buf.WriteString(formatted)
	return buf.String(), nil
}

func (a *Appender) renderNewFileHeader(now time.Time) (string, error) {
	header, err := tmpl.Render(a.templatePath, tmpl.Context{
		Date:      now.Format("Monday, 2 January 2006"),
		CreatedAt: now.Format("2006-01-02 15:04"),
		UpdatedAt: now.Format("2006-01-02 15:04"),
	})
	if err != nil {
		return "", fmt.Errorf("rendering template: %w", err)
	}
	return header, nil
}

func appendFile(path, content string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("opening log file: %w", err)
	}

	_, writeErr := f.WriteString(content)
	closeErr := f.Close()
	if writeErr != nil {
		return fmt.Errorf("writing to log file: %w", writeErr)
	}
	if closeErr != nil {
		return fmt.Errorf("closing log file: %w", closeErr)
	}
	return nil
}

// FormatTime formats a time for subheading display per the configured format.
func FormatTime(t time.Time, format config.TimeFormat) string {
	if format == config.TimeFormat12h {
		return t.Format("3:04 PM")
	}
	return t.Format("15:04")
}
