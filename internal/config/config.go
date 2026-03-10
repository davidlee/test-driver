// Package config loads and validates im configuration.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// EditorTimestamp controls which moment is recorded for editor-composed entries.
type EditorTimestamp string

// EditorTimestamp values.
const (
	EditorTimestampStart EditorTimestamp = "start"
	EditorTimestampEnd   EditorTimestamp = "end"
)

// TimestampRounding controls how entry timestamps are rounded.
type TimestampRounding string

// TimestampRounding values.
const (
	TimestampRoundingAdaptive TimestampRounding = "adaptive"
	TimestampRoundingRound10  TimestampRounding = "round10"
)

// TimeFormat controls how timestamp subheadings are displayed.
type TimeFormat string

// TimeFormat values.
const (
	TimeFormat24h TimeFormat = "24h"
	TimeFormat12h TimeFormat = "12h"
)

// Config holds all im configuration.
type Config struct {
	LogDir            string            `toml:"log_dir"`
	Editor            string            `toml:"editor"`
	EditorTimestamp   EditorTimestamp    `toml:"editor_timestamp"`
	TimestampRounding TimestampRounding  `toml:"timestamp_rounding"`
	TimeFormat        TimeFormat         `toml:"time_format"`
	TitleFormat       string             `toml:"title_format"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		LogDir:            "~/log",
		Editor:            "",
		EditorTimestamp:   EditorTimestampStart,
		TimestampRounding: TimestampRoundingAdaptive,
		TimeFormat:        TimeFormat24h,
		TitleFormat:       "",
	}
}

// DefaultConfigPath returns the standard config file location.
func DefaultConfigPath() string {
	if dir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(dir, "im", "im.toml")
	}
	return filepath.Join("~", ".config", "im", "im.toml")
}

// Load reads configuration from path, falling back to defaults for missing
// keys. If the file does not exist, defaults are returned with no error.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()

	expanded, err := expandHome(path)
	if err != nil {
		return cfg, fmt.Errorf("expanding config path: %w", err)
	}

	data, err := os.ReadFile(expanded)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("reading config: %w", err)
	}

	if err := toml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parsing config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// ResolvedLogDir returns LogDir with ~ expanded to the user's home directory.
func (c Config) ResolvedLogDir() (string, error) {
	return expandHome(c.LogDir)
}

// ResolvedEditor returns the editor command: config value, then $EDITOR,
// then $VISUAL. Returns empty string if none are set.
func (c Config) ResolvedEditor() string {
	if c.Editor != "" {
		return c.Editor
	}
	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}
	return os.Getenv("VISUAL")
}

// DefaultTemplatePath returns the standard template file location.
func DefaultTemplatePath() string {
	if dir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(dir, "im", "template.md")
	}
	return filepath.Join("~", ".config", "im", "template.md")
}

func (c Config) validate() error {
	switch c.EditorTimestamp {
	case EditorTimestampStart, EditorTimestampEnd:
	default:
		return fmt.Errorf("invalid editor_timestamp: %q (want %q or %q)",
			c.EditorTimestamp, EditorTimestampStart, EditorTimestampEnd)
	}

	switch c.TimestampRounding {
	case TimestampRoundingAdaptive, TimestampRoundingRound10:
	default:
		return fmt.Errorf("invalid timestamp_rounding: %q (want %q or %q)",
			c.TimestampRounding, TimestampRoundingAdaptive, TimestampRoundingRound10)
	}

	switch c.TimeFormat {
	case TimeFormat24h, TimeFormat12h:
	default:
		return fmt.Errorf("invalid time_format: %q (want %q or %q)",
			c.TimeFormat, TimeFormat24h, TimeFormat12h)
	}

	return nil
}

func expandHome(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolving home directory: %w", err)
	}

	return filepath.Join(home, path[1:]), nil
}
