// Package main is the entrypoint for the im CLI.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	imcli "im/internal/cli"
	"im/internal/config"
	"im/internal/editor"
	"im/internal/logfile"
	"im/internal/reader"

	"github.com/urfave/cli/v3"
)

var version = "dev"

func main() {
	cfg, err := config.Load(config.DefaultConfigPath())
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	app := &cli.Command{
		Name:    "im",
		Usage:   "Interstitial mindfulness logger",
		Version: version,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "task",
				Aliases: []string{"t"},
				Usage:   "Prefix each line with a Markdown task checkbox",
			},
			&cli.BoolFlag{
				Name:    "read",
				Aliases: []string{"r"},
				Usage:   "Render today's log to stdout",
			},
			&cli.BoolFlag{
				Name:    "pager",
				Aliases: []string{"p"},
				Usage:   "Display today's log in a pager",
			},
			&cli.BoolFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "Print today's log file path",
			},
			&cli.BoolFlag{
				Name:    "edit",
				Aliases: []string{"e"},
				Usage:   "Edit today's log file in $EDITOR",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return run(ctx, cmd, cfg)
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(_ context.Context, cmd *cli.Command, cfg config.Config) error {
	logDir, err := cfg.ResolvedLogDir()
	if err != nil {
		return fmt.Errorf("resolving log dir: %w", err)
	}

	switch {
	case cmd.Bool("file"):
		return runFile(logDir)
	case cmd.Bool("read"):
		return runRead(logDir)
	case cmd.Bool("pager"):
		return runPager(logDir)
	case cmd.Bool("edit"):
		return runEdit(cfg, logDir)
	default:
		return runWrite(cmd, cfg, logDir)
	}
}

func todayPath(logDir string) string {
	return filepath.Join(logDir, time.Now().Format("2006-01-02")+".md")
}

func runFile(logDir string) error {
	fmt.Println(todayPath(logDir))
	return nil
}

func runRead(logDir string) error {
	err := reader.Render(todayPath(logDir))
	if errors.Is(err, reader.ErrNoFile) {
		fmt.Println("no entries for today")
		return nil
	}
	return err
}

func runPager(logDir string) error {
	err := reader.View(todayPath(logDir))
	if errors.Is(err, reader.ErrNoFile) {
		fmt.Println("no entries for today")
		return nil
	}
	return err
}

func runEdit(cfg config.Config, logDir string) error {
	path := todayPath(logDir)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("no entries for today")
		return nil
	}

	return editor.EditFile(path, cfg.Editor, editor.DefaultCheck)
}

func runWrite(cmd *cli.Command, cfg config.Config, logDir string) error {
	args := cmd.Args().Slice()
	mode := imcli.DetectInputMode(args)

	var body string
	switch mode {
	case imcli.ModeInline:
		body = strings.Join(args, " ")
	case imcli.ModePipe:
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("reading stdin: %w", err)
		}
		body = string(data)
	case imcli.ModeEditor:
		return runEditorMode(cfg, logDir, cmd.Bool("task"), editor.Edit, time.Now)
	}

	// Empty input guard.
	if strings.TrimSpace(body) == "" {
		return nil
	}

	appender := logfile.NewAppender(time.Now, cfg)
	return appender.Append(logDir, body, cmd.Bool("task"))
}

// editFunc is the signature for editor.Edit, injectable for testing.
type editFunc func(config.Config, editor.CommandChecker) (string, error)

// runEditorMode opens an external editor, captures the entry, and appends it
// to the daily log file. Timestamp is selected per cfg.EditorTimestamp.
func runEditorMode(cfg config.Config, logDir string, task bool, edit editFunc, clock func() time.Time) error {
	startTime := clock()

	body, err := edit(cfg, editor.DefaultCheck)
	if err != nil {
		return fmt.Errorf("editor: %w", err)
	}

	if strings.TrimSpace(body) == "" {
		return nil
	}

	endTime := clock()

	ts := startTime
	if cfg.EditorTimestamp == config.EditorTimestampEnd {
		ts = endTime
	}

	appender := logfile.NewAppender(func() time.Time { return ts }, cfg)
	return appender.Append(logDir, body, task)
}
