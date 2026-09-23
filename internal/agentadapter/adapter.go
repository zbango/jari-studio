package agentadapter

import (
	"context"
	"errors"
	"io"
	"os/exec"
	"sync"
)

var ErrInvalidRequest = errors.New("invalid agent session request")

type Manifest struct {
	ID           string   `json:"id"`
	Version      string   `json:"version"`
	Transport    string   `json:"transport"`
	Capabilities []string `json:"capabilities"`
}

type SessionRequest struct {
	WorkingDirectory string
	Prompt           string
}

type Event struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type Session interface {
	Events() <-chan Event
	WriteInput([]byte) error
	Resize(cols, rows uint16) error
	Wait() error
	Interrupt() error
	Close() error
}

type Adapter interface {
	Manifest() Manifest
	Start(context.Context, SessionRequest) (Session, error)
}

type CLIAdapter struct {
	manifest  Manifest
	command   string
	buildArgs func(string) []string
}

func NewCLIAdapter(manifest Manifest, command string, buildArgs func(string) []string) (*CLIAdapter, error) {
	if manifest.ID == "" || command == "" || buildArgs == nil {
		return nil, ErrInvalidRequest
	}
	return &CLIAdapter{manifest: manifest, command: command, buildArgs: buildArgs}, nil
}

func (a *CLIAdapter) Manifest() Manifest {
	return a.manifest
}

func (a *CLIAdapter) Start(ctx context.Context, request SessionRequest) (Session, error) {
	if request.WorkingDirectory == "" || request.Prompt == "" {
		return nil, ErrInvalidRequest
	}
	command := exec.CommandContext(ctx, a.command, a.buildArgs(request.Prompt)...)
	command.Dir = request.WorkingDirectory
	stdout, err := command.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stdin, err := command.StdinPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err := command.Start(); err != nil {
		return nil, err
	}
	session := &processSession{command: command, stdin: stdin, events: make(chan Event, 32)}
	var readers sync.WaitGroup
	readers.Add(2)
	go session.readOutput(&readers, "stdout", stdout)
	go session.readOutput(&readers, "stderr", stderr)
	go func() {
		readers.Wait()
		session.waitErr = command.Wait()
		close(session.events)
	}()
	return session, nil
}

type processSession struct {
	command *exec.Cmd
	stdin   io.WriteCloser
	events  chan Event
	waitErr error
}

func (s *processSession) Events() <-chan Event { return s.events }

func (s *processSession) Wait() error {
	for range s.events {
	}
	return s.waitErr
}

func (s *processSession) WriteInput(input []byte) error {
	_, err := s.stdin.Write(input)
	return err
}

func (s *processSession) Resize(_, _ uint16) error {
	return errors.New("resize is only available for PTY sessions")
}

func (s *processSession) Interrupt() error {
	if s.command.Process == nil {
		return errors.New("agent process has not started")
	}
	return s.command.Process.Signal(interruptSignal())
}

func (s *processSession) Close() error {
	if s.command.Process == nil {
		return nil
	}
	return s.command.Process.Kill()
}

func (s *processSession) readOutput(readers *sync.WaitGroup, eventType string, reader io.Reader) {
	defer readers.Done()
	buffer := make([]byte, 4096)
	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			s.events <- Event{Type: eventType, Data: string(buffer[:n])}
		}
		if err != nil {
			return
		}
	}
}
