package cli_test

import (
	"os"
	"testing"

	"im/internal/cli"
)

func TestDetectInputModeInline(t *testing.T) {
	t.Parallel()

	got := cli.DetectInputMode([]string{"hello", "world"})
	if got != cli.ModeInline {
		t.Errorf("got %v, want ModeInline", got)
	}
}

func TestDetectInputModeSingleArg(t *testing.T) {
	t.Parallel()

	got := cli.DetectInputMode([]string{"note"})
	if got != cli.ModeInline {
		t.Errorf("got %v, want ModeInline", got)
	}
}

func TestDetectInputModePipe(t *testing.T) {
	t.Parallel()

	// Stub isTerminal to return false (simulating piped stdin).
	cleanup := stubIsTerminal(false)
	defer cleanup()

	got := cli.DetectInputMode(nil)
	if got != cli.ModePipe {
		t.Errorf("got %v, want ModePipe", got)
	}
}

func TestDetectInputModeEditor(t *testing.T) {
	t.Parallel()

	// Stub isTerminal to return true (simulating interactive TTY).
	cleanup := stubIsTerminal(true)
	defer cleanup()

	got := cli.DetectInputMode(nil)
	if got != cli.ModeEditor {
		t.Errorf("got %v, want ModeEditor", got)
	}
}

func TestDetectInputModeArgsWinOverPipe(t *testing.T) {
	t.Parallel()

	// Even with non-TTY stdin, args take priority.
	cleanup := stubIsTerminal(false)
	defer cleanup()

	got := cli.DetectInputMode([]string{"from", "args"})
	if got != cli.ModeInline {
		t.Errorf("got %v, want ModeInline (args should win over pipe)", got)
	}
}

func TestInputModeString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		mode cli.InputMode
		want string
	}{
		{cli.ModeInline, "inline"},
		{cli.ModePipe, "pipe"},
		{cli.ModeEditor, "editor"},
		{cli.InputMode(99), "InputMode(99)"},
	}

	for _, tt := range tests {
		if got := tt.mode.String(); got != tt.want {
			t.Errorf("InputMode(%d).String() = %q, want %q", int(tt.mode), got, tt.want)
		}
	}
}

// stubIsTerminal replaces the isTerminal function for testing and returns
// a cleanup function that restores the original.
func stubIsTerminal(result bool) func() {
	orig := cli.IsTerminal
	cli.IsTerminal = func(_ *os.File) bool { return result }

	return func() { cli.IsTerminal = orig }
}
