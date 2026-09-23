//go:build windows

package terminal

import "context"

func (Manager) Start(context.Context, string, []string, string, []string, uint16, uint16) (Session, error) {
	return nil, ErrUnsupported
}
