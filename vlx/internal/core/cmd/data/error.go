package data

import "fmt"

// CmdError wraps an exec.ExitError with command context.
type CmdError struct {
	Name     string
	Args     []string
	ExitCode int
	Stderr   []byte
	err      error
}

// Error implements the error interface.
func (e *CmdError) Error() string {
	if len(e.Stderr) > 0 {
		return fmt.Sprintf("%s %v: exit code %d: %s", e.Name, e.Args, e.ExitCode, string(e.Stderr))
	}
	return fmt.Sprintf("%s %v: exit code %d", e.Name, e.Args, e.ExitCode)
}

// Unwrap returns the underlying exec.ExitError.
func (e *CmdError) Unwrap() error {
	return e.err
}

// NewCmdError creates a new CmdError.
func NewCmdError(name string, args []string, exitCode int, stderr []byte, err error) *CmdError {
	return &CmdError{
		Name:     name,
		Args:     args,
		ExitCode: exitCode,
		Stderr:   stderr,
		err:      err,
	}
}
