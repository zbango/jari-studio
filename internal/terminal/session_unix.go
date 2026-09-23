//go:build !windows

package terminal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/creack/pty/v2"
)

type unixSession struct {
	cmd       *exec.Cmd
	pty       *os.File
	events    chan Event
	closeOnce sync.Once
	waitOnce  sync.Once
	waitErr   error
}

func (Manager) Start(ctx context.Context, command string, args []string, directory string, environment []string, cols, rows uint16) (Session, error) {
	if command == "" || directory == "" {
		return nil, errors.New("command and working directory are required")
	}
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Dir = directory
	cmd.Env = append(os.Environ(), environment...)
	terminalSize := &pty.Winsize{Cols: cols, Rows: rows}
	if cols == 0 {
		terminalSize.Cols = 120
	}
	if rows == 0 {
		terminalSize.Rows = 40
	}
	file, err := pty.StartWithSize(cmd, terminalSize)
	if err != nil {
		var pathErr *os.PathError
		if errors.As(err, &pathErr) && errors.Is(pathErr.Err, os.ErrNotExist) {
			return nil, fmt.Errorf("executable %q was not found or could not be opened; reinstall it or choose another command: %w", command, err)
		}
		return nil, err
	}
	session := &unixSession{cmd: cmd, pty: file, events: make(chan Event, 64)}
	go session.readOutput()
	return session, nil
}

func (s *unixSession) Events() <-chan Event { return s.events }

func (s *unixSession) WriteInput(input []byte) error {
	_, err := s.pty.Write(input)
	return err
}

func (s *unixSession) Resize(cols, rows uint16) error {
	return pty.Setsize(s.pty, &pty.Winsize{Cols: cols, Rows: rows})
}

func (s *unixSession) Interrupt() error {
	if s.cmd.Process == nil {
		return errors.New("terminal process has not started")
	}
	return s.cmd.Process.Signal(os.Interrupt)
}

func (s *unixSession) Wait() error {
	s.waitOnce.Do(func() {
		s.waitErr = s.cmd.Wait()
		s.close()
	})
	return s.waitErr
}

func (s *unixSession) Close() error {
	if s.cmd.Process != nil {
		_ = s.cmd.Process.Kill()
	}
	s.close()
	return nil
}

func (s *unixSession) close() {
	s.closeOnce.Do(func() {
		_ = s.pty.Close()
		close(s.events)
	})
}

func (s *unixSession) readOutput() {
	defer func() {
		_ = s.Wait()
	}()
	buffer := make([]byte, 4096)
	for {
		count, err := s.pty.Read(buffer)
		if count > 0 {
			s.events <- Event{Type: "output", Data: string(buffer[:count])}
		}
		if err != nil {
			if !errors.Is(err, io.EOF) {
				s.events <- Event{Type: "error", Data: err.Error()}
			}
			return
		}
	}
}
