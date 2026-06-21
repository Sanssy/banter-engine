package logging

import (
	"bytes"
	"strings"
	"testing"
)

func TestLoggerWritesStructuredFields(t *testing.T) {
	var output bytes.Buffer
	logger := New(&output, "scheduler")
	logger.Info("poll %d", 1)

	line := output.String()
	for _, field := range []string{"level=INFO", "component=scheduler", `message="poll 1"`} {
		if !strings.Contains(line, field) {
			t.Fatalf("log line %q does not contain %q", line, field)
		}
	}
}
