package main

import (
	"strings"
	"testing"
)

func TestRunRejectsUnknownCommand(t *testing.T) {
	err := run([]string{"unknown"})
	if err == nil || !strings.Contains(err.Error(), "usage") {
		t.Fatalf("run() error = %v, want usage error", err)
	}
}
