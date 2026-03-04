// Package cli defines the im command-line interface and input mode detection.
package cli

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// InputMode describes how the user is providing entry text.
type InputMode int

// InputMode values.
const (
	ModeInline InputMode = iota // positional args or after --
	ModePipe                    // stdin is piped
	ModeEditor                  // interactive editor session
)

// String returns a human-readable label for the mode.
func (m InputMode) String() string {
	switch m {
	case ModeInline:
		return "inline"
	case ModePipe:
		return "pipe"
	case ModeEditor:
		return "editor"
	default:
		return fmt.Sprintf("InputMode(%d)", int(m))
	}
}

// DetectInputMode determines how entry text will be sourced.
// Priority: explicit args > piped stdin > editor.
func DetectInputMode(args []string) InputMode {
	if len(args) > 0 {
		return ModeInline
	}

	if !IsTerminal(os.Stdin) {
		return ModePipe
	}

	return ModeEditor
}

// IsTerminal reports whether f is connected to a terminal.
// Exported as a variable for testability.
var IsTerminal = func(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}
