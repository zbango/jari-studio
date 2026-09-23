//go:build windows

package agentadapter

import "os"

func interruptSignal() os.Signal { return os.Kill }
