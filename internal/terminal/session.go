package terminal

import "errors"

var ErrUnsupported = errors.New("real PTY is not supported on this platform yet")

type State string

const (
	StateStarting State = "starting"
	StateRunning  State = "running"
	StateExited   State = "exited"
	StateFailed   State = "failed"
	StateClosed   State = "closed"
)

type Event struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type Info struct {
	ID        string `json:"id"`
	Command   string `json:"command"`
	Directory string `json:"directory"`
	State     State  `json:"state"`
}

type Session interface {
	Events() <-chan Event
	WriteInput([]byte) error
	Resize(cols, rows uint16) error
	Interrupt() error
	Wait() error
	Close() error
}

type Manager struct{}
