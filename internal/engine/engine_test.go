package engine

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Sanssy/banter-engine/internal/config"
	"github.com/Sanssy/banter-engine/internal/logging"
)

func TestRunStopsWhenContextIsCanceled(t *testing.T) {
	var output bytes.Buffer
	runtime := &Engine{
		config: config.Config{PollInterval: time.Minute},
		logger: logging.New(&output, "engine"),
		output: &output,
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := runtime.Run(ctx); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !strings.Contains(output.String(), "scheduler stopped") {
		t.Fatalf("Run() output = %q, want shutdown log", output.String())
	}
}

func TestPublishDryRunWritesToOutput(t *testing.T) {
	var output bytes.Buffer
	runtime := &Engine{
		config: config.Config{DryRun: true},
		output: &output,
	}

	if err := runtime.publish("test message"); err != nil {
		t.Fatalf("publish() error = %v", err)
	}
	if output.String() != "test message\n" {
		t.Fatalf("publish() output = %q", output.String())
	}
}
