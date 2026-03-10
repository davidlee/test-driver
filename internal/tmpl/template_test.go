package tmpl

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRender_Default(t *testing.T) {
	// Non-existent path triggers default template.
	got, err := Render("/nonexistent/template.md", Context{
		Date: "Monday, 2 January 2006",
	})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	want := "# Monday, 2 January 2006\n"
	if got != want {
		t.Errorf("Render() = %q, want %q", got, want)
	}
}

func TestRender_CustomTemplate(t *testing.T) {
	dir := t.TempDir()
	tmplPath := filepath.Join(dir, "template.md")
	if err := os.WriteFile(tmplPath, []byte("---\ncreated: {{ .CreatedAt }}\n---\n# {{ .Date }}\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := Render(tmplPath, Context{
		Date:      "Friday, 9 March 2026",
		CreatedAt: "2026-03-09 14:30",
	})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	want := "---\ncreated: 2026-03-09 14:30\n---\n# Friday, 9 March 2026\n"
	if got != want {
		t.Errorf("Render() = %q, want %q", got, want)
	}
}

func TestRender_InvalidTemplateFallsBack(t *testing.T) {
	dir := t.TempDir()
	tmplPath := filepath.Join(dir, "template.md")
	// Invalid template syntax.
	if err := os.WriteFile(tmplPath, []byte("{{ .Unclosed"), 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := Render(tmplPath, Context{
		Date: "Monday, 2 January 2006",
	})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	// Should fall back to default.
	want := "# Monday, 2 January 2006\n"
	if got != want {
		t.Errorf("Render() = %q, want %q", got, want)
	}
}

func TestRender_MissingKeyProducesEmpty(t *testing.T) {
	dir := t.TempDir()
	tmplPath := filepath.Join(dir, "template.md")
	if err := os.WriteFile(tmplPath, []byte("id: {{ .ID }}\n# {{ .Date }}\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := Render(tmplPath, Context{
		Date: "Friday, 9 March 2026",
		// ID not set — should render as empty.
	})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	want := "id: \n# Friday, 9 March 2026\n"
	if got != want {
		t.Errorf("Render() = %q, want %q", got, want)
	}
}
