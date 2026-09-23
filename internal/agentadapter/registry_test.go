package agentadapter

import "testing"

func TestBuiltInManifestsContainRequiredAdapters(t *testing.T) {
	manifests := BuiltInManifests()
	seen := make(map[string]bool, len(manifests))
	for _, manifest := range manifests {
		seen[manifest.ID] = true
	}
	for _, required := range []string{"codex", "claude-code", "opencode", "generic-pty", "openai-compatible-local"} {
		if !seen[required] {
			t.Errorf("missing built-in adapter %q", required)
		}
	}
}
