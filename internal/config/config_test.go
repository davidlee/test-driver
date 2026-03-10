package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"im/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := config.DefaultConfig()

	if cfg.LogDir != "~/log" {
		t.Errorf("LogDir = %q, want %q", cfg.LogDir, "~/log")
	}
	if cfg.Editor != "" {
		t.Errorf("Editor = %q, want empty", cfg.Editor)
	}
	if cfg.EditorTimestamp != config.EditorTimestampStart {
		t.Errorf("EditorTimestamp = %q, want %q", cfg.EditorTimestamp, config.EditorTimestampStart)
	}
	if cfg.TimestampRounding != config.TimestampRoundingAdaptive {
		t.Errorf("TimestampRounding = %q, want %q", cfg.TimestampRounding, config.TimestampRoundingAdaptive)
	}
	if cfg.TimeFormat != config.TimeFormat24h {
		t.Errorf("TimeFormat = %q, want %q", cfg.TimeFormat, config.TimeFormat24h)
	}
	if cfg.TitleFormat != "" {
		t.Errorf("TitleFormat = %q, want empty", cfg.TitleFormat)
	}
}

func TestLoadMissingFile(t *testing.T) {
	t.Parallel()

	cfg, err := config.Load(filepath.Join(t.TempDir(), "nonexistent.toml"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := config.DefaultConfig()
	if cfg != want {
		t.Errorf("got %+v, want defaults %+v", cfg, want)
	}
}

func TestLoadFullConfig(t *testing.T) {
	t.Parallel()

	content := `
log_dir = "/tmp/mylog"
editor = "nvim"
editor_timestamp = "end"
timestamp_rounding = "round10"
`
	path := writeTempConfig(t, content)

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.LogDir != "/tmp/mylog" {
		t.Errorf("LogDir = %q, want %q", cfg.LogDir, "/tmp/mylog")
	}
	if cfg.Editor != "nvim" {
		t.Errorf("Editor = %q, want %q", cfg.Editor, "nvim")
	}
	if cfg.EditorTimestamp != config.EditorTimestampEnd {
		t.Errorf("EditorTimestamp = %q, want %q", cfg.EditorTimestamp, config.EditorTimestampEnd)
	}
	if cfg.TimestampRounding != config.TimestampRoundingRound10 {
		t.Errorf("TimestampRounding = %q, want %q", cfg.TimestampRounding, config.TimestampRoundingRound10)
	}
}

func TestLoadPartialConfig(t *testing.T) {
	t.Parallel()

	content := `log_dir = "/var/log/im"` + "\n"
	path := writeTempConfig(t, content)

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.LogDir != "/var/log/im" {
		t.Errorf("LogDir = %q, want %q", cfg.LogDir, "/var/log/im")
	}
	// Unset keys should retain defaults.
	if cfg.EditorTimestamp != config.EditorTimestampStart {
		t.Errorf("EditorTimestamp = %q, want default %q", cfg.EditorTimestamp, config.EditorTimestampStart)
	}
}

func TestLoadInvalidEditorTimestamp(t *testing.T) {
	t.Parallel()

	content := `editor_timestamp = "middle"` + "\n"
	path := writeTempConfig(t, content)

	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid editor_timestamp, got nil")
	}
}

func TestLoadInvalidTimestampRounding(t *testing.T) {
	t.Parallel()

	content := `timestamp_rounding = "round5"` + "\n"
	path := writeTempConfig(t, content)

	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid timestamp_rounding, got nil")
	}
}

func TestLoadInvalidTimeFormat(t *testing.T) {
	t.Parallel()

	content := `time_format = "military"` + "\n"
	path := writeTempConfig(t, content)

	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid time_format, got nil")
	}
}

func TestLoadValidTimeFormat12h(t *testing.T) {
	t.Parallel()

	content := `time_format = "12h"` + "\n"
	path := writeTempConfig(t, content)

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.TimeFormat != config.TimeFormat12h {
		t.Errorf("TimeFormat = %q, want %q", cfg.TimeFormat, config.TimeFormat12h)
	}
}

func TestLoadInvalidTOML(t *testing.T) {
	t.Parallel()

	path := writeTempConfig(t, `[[[broken`)

	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid TOML, got nil")
	}
}

func TestResolvedLogDir(t *testing.T) {
	t.Parallel()

	cfg := config.Config{LogDir: "/absolute/path"}
	got, err := cfg.ResolvedLogDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/absolute/path" {
		t.Errorf("got %q, want %q", got, "/absolute/path")
	}
}

func TestResolvedLogDirExpandsHome(t *testing.T) {
	t.Parallel()

	cfg := config.Config{LogDir: "~/log"}
	got, err := cfg.ResolvedLogDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	home, _ := os.UserHomeDir()
	want := filepath.Join(home, "log")
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestResolvedEditorFallback(t *testing.T) {
	// Cannot use t.Parallel — t.Setenv mutates process env.
	cfg := config.Config{Editor: ""}

	t.Setenv("EDITOR", "")
	t.Setenv("VISUAL", "code")

	got := cfg.ResolvedEditor()
	if got != "code" {
		t.Errorf("got %q, want %q", got, "code")
	}
}

func TestResolvedEditorConfigWins(t *testing.T) {
	// Cannot use t.Parallel — t.Setenv mutates process env.
	cfg := config.Config{Editor: "hx"}

	t.Setenv("EDITOR", "nvim")

	got := cfg.ResolvedEditor()
	if got != "hx" {
		t.Errorf("got %q, want %q", got, "hx")
	}
}

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "im.toml")

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}

	return path
}
