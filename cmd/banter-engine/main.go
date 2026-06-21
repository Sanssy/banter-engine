package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/DSanoussy/banter-engine/internal/config"
	"github.com/DSanoussy/banter-engine/internal/engine"
	"github.com/DSanoussy/banter-engine/internal/logging"
)

func main() {
	logger := logging.New(os.Stderr, "main")
	if err := run(os.Args[1:]); err != nil {
		logger.Error("%v", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) != 1 || (args[0] != "run" && args[0] != "dry-run") {
		return fmt.Errorf("usage: banter-engine <run|dry-run>")
	}

	if args[0] == "dry-run" {
		if err := os.Setenv("DRY_RUN", "true"); err != nil {
			return fmt.Errorf("enable dry run: %w", err)
		}
	}
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	runtime, err := engine.New(cfg, os.Stdout)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return runtime.Run(ctx)
}
