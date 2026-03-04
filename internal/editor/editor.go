// Package editor provides external editor integration for entry composition.
package editor

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"

	goeditor "github.com/confluentinc/go-editor"

	"im/internal/config"
)

// DefaultEditors is the fallback chain when no config or $EDITOR is set.
var DefaultEditors = []string{
	"nvim", "hx", "helix", "vim", "emacs", "nano", "code --wait",
}

// CommandChecker reports whether a command is available on PATH.
type CommandChecker func(cmd string) bool

// DefaultCheck uses exec.LookPath to check command availability.
func DefaultCheck(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// Edit opens an external editor for entry composition. Returns the edited
// content. Returns ("", nil) when the user aborts or saves empty content.
func Edit(cfg config.Config, check CommandChecker) (string, error) {
	editorCmd, err := detectEditor(cfg.Editor, check)
	if err != nil {
		return "", fmt.Errorf("detecting editor: %w", err)
	}

	content, exitCode, err := launchEditor(editorCmd)
	if err != nil {
		return "", err
	}

	if !shouldSave(exitCode, content) {
		return "", nil
	}

	if err := validateContent(content); err != nil {
		return "", fmt.Errorf("content validation: %w", err)
	}

	return content, nil
}

func detectEditor(cfgEditor string, check CommandChecker) (string, error) {
	var candidates []string
	if cfgEditor != "" {
		candidates = append([]string{cfgEditor}, DefaultEditors...)
	} else if env := os.Getenv("EDITOR"); env != "" {
		candidates = append([]string{env}, DefaultEditors...)
	} else {
		candidates = DefaultEditors
	}

	for _, ed := range candidates {
		bin := strings.Fields(ed)[0]
		if check(bin) {
			return ed, nil
		}
	}

	return "", fmt.Errorf("no editor found; set $EDITOR or install one of: %v", DefaultEditors)
}

func launchEditor(editorCmd string) (content string, exitCode int, err error) {
	ed := goeditor.NewEditor()
	ed.Command = editorCmd

	raw, tmpPath, err := ed.LaunchTempFile("im-*.md", strings.NewReader(""))
	if tmpPath != "" {
		defer os.Remove(tmpPath) //nolint:errcheck // best-effort cleanup
	}

	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return "", exitErr.ExitCode(), nil
		}
		return "", 1, fmt.Errorf("editor failed: %w", err)
	}

	return string(raw), 0, nil
}

func shouldSave(exitCode int, content string) bool {
	return exitCode == 0 && strings.TrimSpace(content) != ""
}

func validateContent(content string) error {
	if !utf8.ValidString(content) {
		return fmt.Errorf("content is not valid UTF-8")
	}

	const maxSize = 10 * 1024 * 1024
	if len(content) > maxSize {
		return fmt.Errorf("content too large: %d bytes (max %d)", len(content), maxSize)
	}

	if containsBinary(content) {
		return fmt.Errorf("content appears to contain binary data")
	}

	return nil
}

func containsBinary(content string) bool {
	if strings.Contains(content, "\x00") {
		return true
	}

	nonPrintable := 0
	for _, r := range content {
		if r < 32 && r != '\t' && r != '\n' && r != '\r' {
			nonPrintable++
		}
	}

	return content != "" && float64(nonPrintable)/float64(len(content)) > 0.05
}
