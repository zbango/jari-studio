package agentadapter

import "fmt"

func BuiltInManifests() []Manifest {
	return []Manifest{
		{ID: "codex", Version: "v2", Transport: "jsonl", Capabilities: []string{"structured_events", "working_directory", "diff_events", "changed_files", "interrupt", "resume", "text_input", "file_input", "image_input"}},
		{ID: "claude-code", Version: "v1", Transport: "stream-json", Capabilities: []string{"structured_events", "working_directory", "changed_files", "interrupt", "resume", "text_input", "file_input", "image_input"}},
		{ID: "opencode", Version: "v1", Transport: "acp", Capabilities: []string{"structured_events", "working_directory", "diff_events", "changed_files", "interrupt", "resume", "text_input", "file_input", "local_model_provider"}},
		{ID: "generic-pty", Version: "v1", Transport: "pty", Capabilities: []string{"interactive", "working_directory", "text_input"}},
		{ID: "openai-compatible-local", Version: "v1", Transport: "jsonl", Capabilities: []string{"structured_events", "working_directory", "diff_events", "changed_files", "interrupt", "resume", "text_input", "file_input", "local_model_provider"}},
	}
}

func NewBuiltInAdapter(id string) (Adapter, error) {
	var manifest Manifest
	var command string
	var buildArgs func(string) []string
	switch id {
	case "codex":
		manifest = BuiltInManifests()[0]
		command = "codex"
		buildArgs = func(prompt string) []string { return []string{"exec", "--json", prompt} }
	case "claude-code":
		manifest = BuiltInManifests()[1]
		command = "claude"
		buildArgs = func(prompt string) []string { return []string{"-p", prompt, "--output-format", "stream-json"} }
	case "opencode":
		manifest = BuiltInManifests()[2]
		command = "opencode"
		buildArgs = func(prompt string) []string { return []string{"run", prompt} }
	default:
		return nil, fmt.Errorf("adapter %q requires explicit local configuration", id)
	}
	return NewCLIAdapter(manifest, command, buildArgs)
}
