package cmd

import (
	"errors"
	"os"
	"os/exec"
)

// Process is a handle to a running command.
type Process struct {
	Cmd *exec.Cmd
}

// Pid returns the process ID, or -1 if not started.
func (p *Process) Pid() int {
	if p.Cmd == nil || p.Cmd.Process == nil {
		return -1
	}
	return p.Cmd.Process.Pid
}

// Wait blocks until the process exits.
func (p *Process) Wait() error {
	if p.Cmd == nil {
		return errors.New("process not started")
	}
	return p.Cmd.Wait()
}

// Kill forcefully terminates the process.
func (p *Process) Kill() error {
	if p.Cmd == nil || p.Cmd.Process == nil {
		return errors.New("process not started")
	}
	return p.Cmd.Process.Kill()
}

// Signal sends a signal to the process.
func (p *Process) Signal(sig os.Signal) error {
	if p.Cmd == nil || p.Cmd.Process == nil {
		return errors.New("process not started")
	}
	return p.Cmd.Process.Signal(sig)
}

// Release detaches from the process, allowing it to outlive the caller.
func (p *Process) Release() error {
	if p.Cmd == nil || p.Cmd.Process == nil {
		return errors.New("process not started")
	}
	return p.Cmd.Process.Release()
}
