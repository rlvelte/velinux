package progress

import "github.com/rlvelte/velinux/vlx/internal/core/guard"

// QuickshellProgress is a backend that uses quickshell progress indicator.
type QuickshellProgress struct{}

func (q *QuickshellProgress) Available() bool {
	return guard.Binaries("quickshell") == nil
}

func (q *QuickshellProgress) Stream(message string) {
	panic("implement me")
}
