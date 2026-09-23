package agentadapter

import (
	"context"
	"strings"
	"testing"
)

func TestCLIAdapterRequiresStructuredManifest(t *testing.T) {
	_, err := NewCLIAdapter(Manifest{}, "codex", func(prompt string) []string { return []string{prompt} })
	if err != ErrInvalidRequest {
		t.Fatalf("expected invalid request, got %v", err)
	}
}

func TestCLIAdapterExposesManifest(t *testing.T) {
	adapter, err := NewCLIAdapter(Manifest{ID: "generic-pty", Version: "1", Transport: "pty", Capabilities: []string{"interactive"}}, "codex", func(prompt string) []string { return []string{prompt} })
	if err != nil {
		t.Fatalf("create adapter: %v", err)
	}
	if adapter.Manifest().ID != "generic-pty" || adapter.Manifest().Capabilities[0] != "interactive" {
		t.Fatalf("unexpected manifest: %+v", adapter.Manifest())
	}
}

func TestPTYAdapterStreamsTerminalOutput(t *testing.T) {
	adapter, err := NewPTYAdapter(Manifest{ID: "test-pty", Transport: "pty", Capabilities: []string{"interactive"}}, "/bin/sh", func(string) []string {
		return []string{"-c", "printf pty-ready"}
	})
	if err != nil {
		t.Fatalf("create pty adapter: %v", err)
	}
	session, err := adapter.Start(context.Background(), SessionRequest{WorkingDirectory: ".", Prompt: "run"})
	if err != nil {
		t.Fatalf("start pty adapter: %v", err)
	}
	var output strings.Builder
	for event := range session.Events() {
		output.WriteString(event.Data)
	}
	if err := session.Wait(); err != nil {
		t.Fatalf("wait pty adapter: %v", err)
	}
	if !strings.Contains(output.String(), "pty-ready") {
		t.Fatalf("output %q does not contain marker", output.String())
	}
}
