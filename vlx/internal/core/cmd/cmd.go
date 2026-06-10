// Package cmd provides a fluent, zero-allocation command execution engine.
package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rvelte/vlx/internal/core/cmd/data"
)

// Cmd describes a single external command invocation.
type Cmd struct {
	name    string
	args    []string
	stdout  io.Writer
	stderr  io.Writer
	stdin   io.Reader
	env     []string // nil = inherit os.Environ
	dir     string
	timeout time.Duration
}

// New creates a new command builder.
func New(name string, args ...string) *Cmd {
	return &Cmd{
		name: name,
		args: args,
	}
}

// Raw returns a human-readable representation of the command.
func (c *Cmd) Raw() string {
	return fmt.Sprintf("%s %v", c.name, c.args)
}

// Args appends data arguments.
func (c *Cmd) Args(args ...string) *Cmd {
	c.args = append(c.args, args...)
	return c
}

// Stdout sets the stdout writer.
func (c *Cmd) Stdout(w io.Writer) *Cmd {
	c.stdout = w
	return c
}

// Stderr sets the stderr writer.
func (c *Cmd) Stderr(w io.Writer) *Cmd {
	c.stderr = w
	return c
}

// Stdin sets the stdin reader.
func (c *Cmd) Stdin(r io.Reader) *Cmd {
	c.stdin = r
	return c
}

// Env adds or overrides a single environment variable.
func (c *Cmd) Env(key, val string) *Cmd {
	if c.env == nil {
		c.env = os.Environ()
	}

	prefix := key + "="
	for i, e := range c.env {
		if strings.HasPrefix(e, prefix) {
			c.env[i] = key + "=" + val
			return c
		}
	}

	c.env = append(c.env, key+"="+val)
	return c
}

// SetEnv replaces the entire environment.
func (c *Cmd) SetEnv(env []string) *Cmd {
	c.env = env
	return c
}

// Dir sets the working directory.
func (c *Cmd) Dir(d string) *Cmd {
	c.dir = d
	return c
}

// Timeout sets a maximum execution duration.
func (c *Cmd) Timeout(d time.Duration) *Cmd {
	c.timeout = d
	return c
}

// Run executes the command and waits for it to finish.
func (c *Cmd) Run(ctx context.Context) error {
	_, err := c.run(ctx, false)
	return err
}

// RunCaptured executes the command and captures its stdout and stderr.
func (c *Cmd) RunCaptured(ctx context.Context) (*data.CmdResult, error) {
	return c.run(ctx, true)
}

// RunParallel runs the command in a background goroutine.
func (c *Cmd) RunParallel() {
	go func() { _ = c.Run(context.Background()) }()
}

// RunProcess launches the command without waiting.
func (c *Cmd) RunProcess(ctx context.Context) (*Process, error) {
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, c.name, c.args...)
	cmd.Dir = c.dir
	cmd.Stdout = c.stdout
	cmd.Stderr = c.stderr
	cmd.Stdin = c.stdin

	if c.env != nil {
		cmd.Env = c.env
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Process{Cmd: cmd}, nil
}

// run is the shared execution engine.
func (c *Cmd) run(ctx context.Context, capture bool) (*data.CmdResult, error) {
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	var stdoutBuf, stderrBuf *bytes.Buffer
	stdoutW := c.stdout
	stderrW := c.stderr

	if capture {
		stdoutBuf = &bytes.Buffer{}
		stderrBuf = &bytes.Buffer{}

		if stdoutW != nil {
			stdoutW = io.MultiWriter(stdoutBuf, stdoutW)
		} else {
			stdoutW = stdoutBuf
		}

		if stderrW != nil {
			stderrW = io.MultiWriter(stderrBuf, stderrW)
		} else {
			stderrW = stderrBuf
		}
	}

	if stdoutW == nil {
		stdoutW = io.Discard
	}

	if stderrW == nil {
		stderrW = io.Discard
	}

	cmd := exec.CommandContext(ctx, c.name, c.args...)
	cmd.Dir = c.dir
	cmd.Stdout = stdoutW
	cmd.Stderr = stderrW
	cmd.Stdin = c.stdin

	if c.env != nil {
		cmd.Env = c.env
	}

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	if err == nil {
		result := &data.CmdResult{
			ExitCode: 0,
			Duration: duration,
		}

		if capture {
			result.Stdout = stdoutBuf.Bytes()
			result.Stderr = stderrBuf.Bytes()
		}

		return result, nil
	}

	if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
		result := &data.CmdResult{
			ExitCode: exitErr.ExitCode(),
			Duration: duration,
		}

		if capture {
			result.Stderr = stderrBuf.Bytes()
		}

		var stderr []byte
		if stderrBuf != nil {
			stderr = stderrBuf.Bytes()
		}
		return result, data.NewCmdError(c.name, c.args, exitErr.ExitCode(), stderr, err)
	}

	return nil, err
}
