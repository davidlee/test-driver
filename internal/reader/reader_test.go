package reader_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"im/internal/reader"
)

func TestResolveViewer_FindsCat(t *testing.T) {
	// cat should always be available on POSIX systems.
	tmp := filepath.Join(t.TempDir(), "test.md")
	if err := os.WriteFile(tmp, []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Unset PAGER to ensure we fall through to cat.
	t.Setenv("PAGER", "")

	bin, argv, err := reader.ResolveViewer(tmp)
	if err != nil {
		t.Fatalf("ResolveViewer() error: %v", err)
	}

	if bin == "" {
		t.Fatal("expected non-empty bin path")
	}

	// Last element of argv should be the file path.
	if len(argv) == 0 || argv[len(argv)-1] != tmp {
		t.Errorf("expected argv to end with file path, got %v", argv)
	}
}

func TestResolveViewer_PagerEnv(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "test.md")
	if err := os.WriteFile(tmp, []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Set PAGER to cat (known to exist) to test env var pickup.
	t.Setenv("PAGER", "cat")

	_, argv, err := reader.ResolveViewer(tmp)
	if err != nil {
		t.Fatalf("ResolveViewer() error: %v", err)
	}

	// argv[0] should be "cat" (from PAGER) or "glow"/"rich" if present.
	// We just verify it resolves without error.
	if len(argv) == 0 {
		t.Fatal("expected non-empty argv")
	}
}

func TestView_NoFile(t *testing.T) {
	err := reader.View("/nonexistent/path/to/file.md")
	if !errors.Is(err, reader.ErrNoFile) {
		t.Errorf("View() error = %v, want ErrNoFile", err)
	}
}

func TestResolveRenderer_FindsCat(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "test.md")
	if err := os.WriteFile(tmp, []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}

	bin, argv, err := reader.ResolveRenderer(tmp)
	if err != nil {
		t.Fatalf("ResolveRenderer() error: %v", err)
	}

	if bin == "" {
		t.Fatal("expected non-empty bin path")
	}

	if len(argv) == 0 || argv[len(argv)-1] != tmp {
		t.Errorf("expected argv to end with file path, got %v", argv)
	}

	// ResolveRenderer must NOT include --pager in argv.
	for _, arg := range argv {
		if arg == "--pager" {
			t.Error("ResolveRenderer must not include --pager")
		}
	}
}

func TestRender_NoFile(t *testing.T) {
	err := reader.Render("/nonexistent/path/to/file.md")
	if !errors.Is(err, reader.ErrNoFile) {
		t.Errorf("Render() error = %v, want ErrNoFile", err)
	}
}
