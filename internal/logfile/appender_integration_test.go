package logfile_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"im/internal/config"
	"im/internal/logfile"
)

// Integration tests — exercise the full write path through Appender.

func TestIntegration_FullWriteSequence(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Config{TimestampRounding: config.TimestampRoundingAdaptive}
	date := time.Date(2026, 1, 23, 0, 0, 0, 0, time.UTC)

	// 10:03 — first entry, new file
	a1 := logfile.NewAppender(func() time.Time {
		return date.Add(10*time.Hour + 3*time.Minute)
	}, cfg)
	if err := a1.Append(dir, "fixed the auth bug", false); err != nil {
		t.Fatalf("entry 1: %v", err)
	}

	// 10:07 — within 10m, not on boundary, suppress heading
	a2 := logfile.NewAppender(func() time.Time {
		return date.Add(10*time.Hour + 7*time.Minute)
	}, cfg)
	if err := a2.Append(dir, "also updated the docs", false); err != nil {
		t.Fatalf("entry 2: %v", err)
	}

	// 10:10 — within 10m but on round boundary, new heading
	a3 := logfile.NewAppender(func() time.Time {
		return date.Add(10*time.Hour + 10*time.Minute)
	}, cfg)
	if err := a3.Append(dir, "standup notes", false); err != nil {
		t.Fatalf("entry 3: %v", err)
	}

	// 10:47 — gap > 10m, new heading with exact time
	a4 := logfile.NewAppender(func() time.Time {
		return date.Add(10*time.Hour + 47*time.Minute)
	}, cfg)
	if err := a4.Append(dir, "reviewed PR #42", false); err != nil {
		t.Fatalf("entry 4: %v", err)
	}

	// 14:30 — task entry
	a5 := logfile.NewAppender(func() time.Time {
		return date.Add(14*time.Hour + 30*time.Minute)
	}, cfg)
	if err := a5.Append(dir, "buy milk\nwalk dog", true); err != nil {
		t.Fatalf("entry 5: %v", err)
	}

	// Empty input — should not modify file
	a6 := logfile.NewAppender(func() time.Time {
		return date.Add(15 * time.Hour)
	}, cfg)
	if err := a6.Append(dir, "   ", false); err != nil {
		t.Fatalf("entry 6: %v", err)
	}

	path := filepath.Join(dir, "2026-01-23.md")
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}

	want := `# Friday, 23 January 2026

## 10:03

fixed the auth bug

also updated the docs

## 10:10

standup notes

## 10:47

reviewed PR #42

## 14:30

- [ ] buy milk
- [ ] walk dog
`
	if string(got) != want {
		t.Errorf("file content mismatch:\n--- got ---\n%s\n--- want ---\n%s", string(got), want)
	}
}

func TestIntegration_Round10FullSequence(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Config{TimestampRounding: config.TimestampRoundingRound10}
	date := time.Date(2026, 1, 23, 0, 0, 0, 0, time.UTC)

	// 10:03 → rounds to 10:00
	a1 := logfile.NewAppender(func() time.Time {
		return date.Add(10*time.Hour + 3*time.Minute)
	}, cfg)
	if err := a1.Append(dir, "first", false); err != nil {
		t.Fatal(err)
	}

	// 10:07 → rounds to 10:00, same as last → suppressed
	a2 := logfile.NewAppender(func() time.Time {
		return date.Add(10*time.Hour + 7*time.Minute)
	}, cfg)
	if err := a2.Append(dir, "second", false); err != nil {
		t.Fatal(err)
	}

	// 10:13 → rounds to 10:10, new heading
	a3 := logfile.NewAppender(func() time.Time {
		return date.Add(10*time.Hour + 13*time.Minute)
	}, cfg)
	if err := a3.Append(dir, "third", false); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(dir, "2026-01-23.md")
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	want := `# Friday, 23 January 2026

## 10:00

first

second

## 10:10

third
`
	if string(got) != want {
		t.Errorf("file content mismatch:\n--- got ---\n%s\n--- want ---\n%s", string(got), want)
	}
}
