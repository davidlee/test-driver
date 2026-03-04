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
				Usage:   "Display today's log in a pager",
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

	if cmd.Bool("read") {
		return runRead(logDir)
	}

	return runWrite(cmd, cfg, logDir)
}

func runRead(logDir string) error {
	today := time.Now().Format("2006-01-02") + ".md"
	path := filepath.Join(logDir, today)

	err := reader.View(path)
	if errors.Is(err, reader.ErrNoFile) {
		fmt.Println("no entries for today")
		return nil
	}
	return err
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
		fmt.Println("editor mode: not yet implemented (DE-003)")
		return nil
	}

	// Empty input guard.
	if strings.TrimSpace(body) == "" {
		return nil
	}

	appender := logfile.NewAppender(time.Now, cfg)
	return appender.Append(logDir, body, cmd.Bool("task"))
}
