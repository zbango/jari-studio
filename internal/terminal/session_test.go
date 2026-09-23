package terminal

import (
	"context"
	"strings"
	"testing"
)

func TestUnixSessionStreamsOutput(t *testing.T) {
	if testing.Short() {
		t.Skip("process session test skipped in short mode")
	}
	session, err := (Manager{}).Start(context.Background(), "/bin/sh", []string{"-c", "printf terminal-ready"}, ".", nil, 100, 30)
	if err != nil {
		t.Fatalf("start terminal: %v", err)
	}
	var output strings.Builder
	for event := range session.Events() {
		if event.Type == "output" {
			output.WriteString(event.Data)
		}
	}
	if err := session.Wait(); err != nil {
		t.Fatalf("wait terminal: %v", err)
	}
	if !strings.Contains(output.String(), "terminal-ready") {
		t.Fatalf("output %q does not contain terminal marker", output.String())
	}
}
