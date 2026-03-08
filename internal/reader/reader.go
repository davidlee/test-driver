// Package reader dispatches today's log file to a viewer.
package reader

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// ErrNoFile is returned when the daily log file does not exist.
var ErrNoFile = errors.New("no entries for today")

// viewer describes a candidate viewer command.
type viewer struct {
	name string
	args []string // args to pass before the file path
}

// pagerCandidates in priority order: glow --pager, rich --markdown --pager, $PAGER, cat.
func pagerCandidates() []viewer {
	out := []viewer{
		{name: "glow", args: []string{"--pager"}},
		{name: "rich", args: []string{"--markdown", "--pager"}},
	}

	if pager := os.Getenv("PAGER"); pager != "" {
		out = append(out, viewer{name: pager, args: nil})
	}

	out = append(out, viewer{name: "cat", args: nil})
	return out
}

// rendererCandidates in priority order: glow (stdout, no --pager), cat.
func rendererCandidates() []viewer {
	return []viewer{
		{name: "glow", args: nil},
		{name: "cat", args: nil},
	}
}

func resolve(filePath string, candidates []viewer) (binPath string, argv []string, err error) {
	for _, v := range candidates {
		bin, lookErr := exec.LookPath(v.name)
		if lookErr != nil {
			continue
		}
		argv := append([]string{v.name}, v.args...)
		argv = append(argv, filePath)
		return bin, argv, nil
	}
	return "", nil, fmt.Errorf("no viewer found")
}

// ResolveViewer finds the first available pager viewer. Exported for testability.
func ResolveViewer(filePath string) (binPath string, argv []string, err error) {
	return resolve(filePath, pagerCandidates())
}

// ResolveRenderer finds the first available stdout renderer. Exported for testability.
func ResolveRenderer(filePath string) (binPath string, argv []string, err error) {
	return resolve(filePath, rendererCandidates())
}

func execViewer(path string, resolver func(string) (string, []string, error)) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ErrNoFile
	}

	bin, argv, err := resolver(path)
	if err != nil {
		return err
	}

	return syscall.Exec(bin, argv, os.Environ())
}

// View displays the file at path in a pager.
// Replaces the current process via syscall.Exec.
func View(path string) error {
	return execViewer(path, ResolveViewer)
}

// Render displays the file at path to stdout (no pager).
// Replaces the current process via syscall.Exec.
func Render(path string) error {
	return execViewer(path, ResolveRenderer)
}
