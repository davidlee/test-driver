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

// candidates in priority order: glow --pager, rich --markdown --pager, $PAGER, cat.
func candidates() []viewer {
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

// ResolveViewer finds the first available viewer and returns the full path
// and argv. Exported for testability.
func ResolveViewer(filePath string) (binPath string, argv []string, err error) {
	for _, v := range candidates() {
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

// View displays the file at path using the first available viewer.
// Replaces the current process via syscall.Exec.
func View(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ErrNoFile
	}

	bin, argv, err := ResolveViewer(path)
	if err != nil {
		return err
	}

	return syscall.Exec(bin, argv, os.Environ())
}
