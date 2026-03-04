// Package main is the entrypoint for the im CLI.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	imcli "im/internal/cli"
	"im/internal/config"

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
	if cmd.Bool("read") {
		fmt.Println("read mode: not yet implemented")
		return nil
	}

	// Gather positional args (text after flags, or after --).
	args := cmd.Args().Slice()
	mode := imcli.DetectInputMode(args)

	logDir, err := cfg.ResolvedLogDir()
	if err != nil {
		return fmt.Errorf("resolving log dir: %w", err)
	}

	fmt.Printf("mode:   %s\n", mode)
	fmt.Printf("task:   %v\n", cmd.Bool("task"))
	fmt.Printf("logdir: %s\n", logDir)

	if mode == imcli.ModeInline {
		fmt.Printf("text:   %s\n", strings.Join(args, " "))
	}

	fmt.Println("(entry writing not yet implemented)")

	return nil
}
