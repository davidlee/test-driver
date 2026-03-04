package logfile_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"im/internal/config"
	"im/internal/logfile"
)

func fixedClock(h, m int) func() time.Time {
	return func() time.Time {
		return time.Date(2026, 1, 23, h, m, 0, 0, time.UTC)
	}
}

func TestAppender_NewFile(t *testing.T) {
	dir := t.TempDir()
	a := logfile.NewAppender(fixedClock(10, 3), config.Config{
		TimestampRounding: config.TimestampRoundingAdaptive,
	})

	err := a.Append(dir, "hello world", false)
	if err != nil {
		t.Fatalf("Append() error: %v", err)
	}

	path := filepath.Join(dir, "2026-01-23.md")
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}

	want := "# Friday, 23 January 2026\n\n## 10:03\n\nhello world\n"
	if string(got) != want {
		t.Errorf("file content:\n got: %q\nwant: %q", string(got), want)
	}
}

func TestAppender_AppendToExisting(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Config{TimestampRounding: config.TimestampRoundingAdaptive}

	// First entry at 10:03.
	a1 := logfile.NewAppender(fixedClock(10, 3), cfg)
	if err := a1.Append(dir, "first", false); err != nil {
		t.Fatalf("first Append: %v", err)
	}

	// Second entry at 10:47 (gap > 10m, new heading).
	a2 := logfile.NewAppender(fixedClock(10, 47), cfg)
	if err := a2.Append(dir, "second", false); err != nil {
		t.Fatalf("second Append: %v", err)
	}

	path := filepath.Join(dir, "2026-01-23.md")
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}

	want := "# Friday, 23 January 2026\n\n## 10:03\n\nfirst\n\n## 10:47\n\nsecond\n"
	if string(got) != want {
		t.Errorf("file content:\n got: %q\nwant: %q", string(got), want)
	}
}

func TestAppender_SuppressHeading(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Config{TimestampRounding: config.TimestampRoundingAdaptive}

	// First entry at 10:03.
	a1 := logfile.NewAppender(fixedClock(10, 3), cfg)
	if err := a1.Append(dir, "first", false); err != nil {
		t.Fatalf("first Append: %v", err)
	}

	// Second entry at 10:07 (gap < 10m, no round boundary, suppress heading).
	a2 := logfile.NewAppender(fixedClock(10, 7), cfg)
	if err := a2.Append(dir, "second", false); err != nil {
		t.Fatalf("second Append: %v", err)
	}

	path := filepath.Join(dir, "2026-01-23.md")
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}

	want := "# Friday, 23 January 2026\n\n## 10:03\n\nfirst\n\nsecond\n"
	if string(got) != want {
		t.Errorf("file content:\n got: %q\nwant: %q", string(got), want)
	}
}

func TestAppender_TaskPrefix(t *testing.T) {
	dir := t.TempDir()
	a := logfile.NewAppender(fixedClock(14, 30), config.Config{
		TimestampRounding: config.TimestampRoundingAdaptive,
	})

	err := a.Append(dir, "buy milk", true)
	if err != nil {
		t.Fatalf("Append() error: %v", err)
	}

	path := filepath.Join(dir, "2026-01-23.md")
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}

	want := "# Friday, 23 January 2026\n\n## 14:30\n\n- [ ] buy milk\n"
	if string(got) != want {
		t.Errorf("file content:\n got: %q\nwant: %q", string(got), want)
	}
}

func TestAppender_Round10Strategy(t *testing.T) {
	dir := t.TempDir()
	a := logfile.NewAppender(fixedClock(10, 7), config.Config{
		TimestampRounding: config.TimestampRoundingRound10,
	})

	err := a.Append(dir, "hello", false)
	if err != nil {
		t.Fatalf("Append() error: %v", err)
	}

	path := filepath.Join(dir, "2026-01-23.md")
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}

	want := "# Friday, 23 January 2026\n\n## 10:00\n\nhello\n"
	if string(got) != want {
		t.Errorf("file content:\n got: %q\nwant: %q", string(got), want)
	}
}

func TestAppender_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "log")
	a := logfile.NewAppender(fixedClock(9, 0), config.Config{
		TimestampRounding: config.TimestampRoundingAdaptive,
	})

	err := a.Append(dir, "hello", false)
	if err != nil {
		t.Fatalf("Append() error: %v", err)
	}

	path := filepath.Join(dir, "2026-01-23.md")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}

func TestAppender_PriorContentUnmodified(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Config{TimestampRounding: config.TimestampRoundingAdaptive}

	// Write first entry.
	a1 := logfile.NewAppender(fixedClock(10, 3), cfg)
	if err := a1.Append(dir, "first", false); err != nil {
		t.Fatalf("first Append: %v", err)
	}

	path := filepath.Join(dir, "2026-01-23.md")
	before, _ := os.ReadFile(path)

	// Write second entry.
	a2 := logfile.NewAppender(fixedClock(10, 47), cfg)
	if err := a2.Append(dir, "second", false); err != nil {
		t.Fatalf("second Append: %v", err)
	}

	after, _ := os.ReadFile(path)

	// After must start with before (append-only).
	if len(after) < len(before) || string(after[:len(before)]) != string(before) {
		t.Errorf("prior content was modified:\nbefore: %q\nafter:  %q", string(before), string(after))
	}
}
