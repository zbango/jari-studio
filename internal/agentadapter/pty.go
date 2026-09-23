package agentadapter

import (
	"context"

	"github.com/agentic-app-studio/studio/internal/terminal"
)

type PTYAdapter struct {
	manifest  Manifest
	command   string
	buildArgs func(string) []string
}

func NewPTYAdapter(manifest Manifest, command string, buildArgs func(string) []string) (*PTYAdapter, error) {
	if manifest.ID == "" || command == "" || buildArgs == nil {
		return nil, ErrInvalidRequest
	}
	return &PTYAdapter{manifest: manifest, command: command, buildArgs: buildArgs}, nil
}

func (a *PTYAdapter) Manifest() Manifest { return a.manifest }

func (a *PTYAdapter) Start(ctx context.Context, request SessionRequest) (Session, error) {
	if request.WorkingDirectory == "" || request.Prompt == "" {
		return nil, ErrInvalidRequest
	}
	process, err := (terminal.Manager{}).Start(ctx, a.command, a.buildArgs(request.Prompt), request.WorkingDirectory, nil, 120, 40)
	if err != nil {
		return nil, err
	}
	session := &ptySession{process: process, events: make(chan Event, 64)}
	go func() {
		for event := range process.Events() {
			session.events <- Event{Type: event.Type, Data: event.Data}
		}
		close(session.events)
	}()
	return session, nil
}

type ptySession struct {
	process terminal.Session
	events  chan Event
}

func (s *ptySession) Events() <-chan Event           { return s.events }
func (s *ptySession) WriteInput(input []byte) error  { return s.process.WriteInput(input) }
func (s *ptySession) Resize(cols, rows uint16) error { return s.process.Resize(cols, rows) }
func (s *ptySession) Wait() error                    { return s.process.Wait() }
func (s *ptySession) Interrupt() error               { return s.process.Interrupt() }
func (s *ptySession) Close() error                   { return s.process.Close() }
