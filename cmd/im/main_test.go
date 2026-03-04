package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"im/internal/config"
	"im/internal/editor"
)

// sequentialClock returns a clock that yields each time in order,
// then repeats the last value.
func sequentialClock(times ...time.Time) func() time.Time {
	i := 0
	return func() time.Time {
		if i >= len(times) {
			return times[len(times)-1]
		}
		t := times[i]
		i++
		return t
	}
}

func mockEdit(body string) editFunc {
	return func(_ config.Config, _ editor.CommandChecker) (string, error) {
		return body, nil
	}
}

func TestRunEditorMode_SavesContent(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Config{
		EditorTimestamp:   config.EditorTimestampStart,
		TimestampRounding: config.TimestampRoundingAdaptive,
	}

	start := time.Date(2026, 3, 4, 14, 30, 0, 0, time.UTC)
	end := start.Add(5 * time.Minute)
	clock := sequentialClock(start, end)

	err := runEditorMode(cfg, dir, false, mockEdit("hello from editor\n"), clock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	path := filepath.Join(dir, "2026-03-04.md")
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading log file: %v", err)
	}

	content := string(got)
	if !strings.Contains(content, "hello from editor") {
		t.Errorf("expected content in log file, got:\n%s", content)
	}
	if !strings.Contains(content, "14:30") {
		t.Errorf("expected start timestamp 14:30, got:\n%s", content)
	}
}

func TestRunEditorMode_EmptyAbort(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Config{
		EditorTimestamp:   config.EditorTimestampStart,
		TimestampRounding: config.TimestampRoundingAdaptive,
	}

	start := time.Date(2026, 3, 4, 14, 30, 0, 0, time.UTC)
	clock := sequentialClock(start)

	err := runEditorMode(cfg, dir, false, mockEdit(""), clock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	path := filepath.Join(dir, "2026-03-04.md")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected no log file for empty/aborted content")
	}
}

func TestRunEditorMode_WhitespaceOnlyAbort(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Config{
		EditorTimestamp:   config.EditorTimestampStart,
		TimestampRounding: config.TimestampRoundingAdaptive,
	}

	start := time.Date(2026, 3, 4, 14, 30, 0, 0, time.UTC)
	clock := sequentialClock(start)

	err := runEditorMode(cfg, dir, false, mockEdit("  \n\t  \n"), clock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	path := filepath.Join(dir, "2026-03-04.md")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected no log file for whitespace-only content")
	}
}

func TestRunEditorMode_TimestampEnd(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Config{
		EditorTimestamp:   config.EditorTimestampEnd,
		TimestampRounding: config.TimestampRoundingAdaptive,
	}

	start := time.Date(2026, 3, 4, 14, 30, 0, 0, time.UTC)
	end := time.Date(2026, 3, 4, 14, 45, 0, 0, time.UTC)
	clock := sequentialClock(start, end)

	err := runEditorMode(cfg, dir, false, mockEdit("end timestamp test\n"), clock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	path := filepath.Join(dir, "2026-03-04.md")
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading log file: %v", err)
	}

	content := string(got)
	if !strings.Contains(content, "14:45") {
		t.Errorf("expected end timestamp 14:45, got:\n%s", content)
	}
	if strings.Contains(content, "14:30") {
		t.Errorf("should not contain start timestamp 14:30, got:\n%s", content)
	}
}

func TestRunEditorMode_TaskCheckbox(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Config{
		EditorTimestamp:   config.EditorTimestampStart,
		TimestampRounding: config.TimestampRoundingAdaptive,
	}

	start := time.Date(2026, 3, 4, 10, 0, 0, 0, time.UTC)
	clock := sequentialClock(start, start.Add(time.Minute))

	err := runEditorMode(cfg, dir, true, mockEdit("buy milk\nwalk dog"), clock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	path := filepath.Join(dir, "2026-03-04.md")
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading log file: %v", err)
	}

	content := string(got)
	if !strings.Contains(content, "- [ ] buy milk") {
		t.Errorf("expected task checkbox, got:\n%s", content)
	}
	if !strings.Contains(content, "- [ ] walk dog") {
		t.Errorf("expected task checkbox for second line, got:\n%s", content)
	}
}
