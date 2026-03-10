// Package tmpl loads and renders Go text/templates for daily log files.
package tmpl

import (
	"bytes"
	_ "embed"
	"os"
	"text/template"
)

//go:embed default.md
var defaultTmpl string

// Context holds the variables available to daily file templates.
type Context struct {
	Date      string // e.g. "Monday, 2 January 2006"
	CreatedAt string // e.g. "2006-01-02 15:04"
	UpdatedAt string // e.g. "2006-01-02 15:04"
	ID        string // 7-char nanoID (empty until P02)
	Title     string // rendered from title_format (empty until P02)
}

// Render loads the template from path (falling back to the bundled default)
// and executes it with the given context.
func Render(templatePath string, ctx Context) (string, error) {
	tmplText := defaultTmpl
	if data, err := os.ReadFile(templatePath); err == nil {
		tmplText = string(data)
	}

	t, err := template.New("daily").Option("missingkey=zero").Parse(tmplText)
	if err != nil {
		// Fall back to default on parse error.
		t = template.Must(template.New("daily").Option("missingkey=zero").Parse(defaultTmpl))
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, ctx); err != nil {
		return "", err
	}
	return buf.String(), nil
}
