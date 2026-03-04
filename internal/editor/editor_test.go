package editor

import (
	"testing"
)

func TestDetectEditor(t *testing.T) {
	available := func(cmds ...string) CommandChecker {
		set := make(map[string]bool, len(cmds))
		for _, c := range cmds {
			set[c] = true
		}
		return func(cmd string) bool { return set[cmd] }
	}

	tests := []struct {
		name      string
		editor    string // config editor field
		envEditor string // $EDITOR
		check     CommandChecker
		want      string
		wantErr   bool
	}{
		{
			name:   "config editor takes priority",
			editor: "hx",
			check:  available("hx", "nvim"),
			want:   "hx",
		},
		{
			name:      "EDITOR env used when no config",
			envEditor: "emacs",
			check:     available("emacs", "nvim"),
			want:      "emacs",
		},
		{
			name:  "falls through to defaults",
			check: available("vim"),
			want:  "vim",
		},
		{
			name:  "first available default wins",
			check: available("nvim", "vim"),
			want:  "nvim",
		},
		{
			name:    "error when nothing available",
			check:   available(),
			wantErr: true,
		},
		{
			name:   "config editor not available falls through",
			editor: "zed",
			check:  available("nvim"),
			want:   "nvim",
		},
		{
			name:      "EDITOR not available falls through",
			envEditor: "micro",
			check:     available("vim"),
			want:      "vim",
		},
		{
			name:   "config with flags preserves flags",
			editor: "code --wait",
			check:  available("code"),
			want:   "code --wait",
		},
		{
			name:  "default code --wait resolved",
			check: available("code"),
			want:  "code --wait",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envEditor != "" {
				t.Setenv("EDITOR", tt.envEditor)
			} else {
				t.Setenv("EDITOR", "")
			}

			got, err := detectEditor(tt.editor, tt.check)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestShouldSave(t *testing.T) {
	tests := []struct {
		name     string
		exitCode int
		content  string
		want     bool
	}{
		{"exit 0 with content", 0, "hello\n", true},
		{"exit 0 empty", 0, "", false},
		{"exit 0 whitespace only", 0, "  \n\t  ", false},
		{"exit 1 with content", 1, "hello\n", false},
		{"exit 1 empty", 1, "", false},
		{"exit 130 (sigint)", 130, "partial", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldSave(tt.exitCode, tt.content)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{"valid text", "hello world\n", false},
		{"valid with tabs and newlines", "line1\n\tindented\nline3\n", false},
		{"empty string", "", false},
		{"invalid utf8", "hello \xff world", true},
		{"null bytes", "hello\x00world", true},
		{"high non-printable ratio", string(make([]byte, 100)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateContent(tt.content)
			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateContent_TooLarge(t *testing.T) {
	large := make([]byte, 10*1024*1024+1)
	for i := range large {
		large[i] = 'a'
	}
	if err := validateContent(string(large)); err == nil {
		t.Fatal("expected error for oversized content")
	}
}
