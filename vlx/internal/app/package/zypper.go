package _package

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Zypper struct{}

// Search searches for packages with structured XML output.
func (z *Zypper) Search(ctx context.Context, query string) ([]Package, error) {
	args := []string{"--xmlout", "search", "--details", query}

	cmd := z.command(ctx, true, false, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
			if exitErr.ExitCode() == 2 || exitErr.ExitCode() == 104 {
				return []Package{}, nil
			}
		}

		return nil, fmt.Errorf("zypper search failed: %w", err)
	}

	return parseXml(stdout.Bytes())
}

// Info retrieves detailed information about a package.
func (z *Zypper) Info(ctx context.Context, pkgName string) (string, error) {
	cmd := z.command(ctx, true, false, "info", pkgName)
	out, err := cmd.Output()

	if err != nil {
		if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
			return "", fmt.Errorf("zypper info %s failed: %s", pkgName, strings.TrimSpace(string(exitErr.Stderr)))
		}
		return "", fmt.Errorf("zypper info %s failed: %w", pkgName, err)
	}

	return strings.TrimLeft(string(out), " \t\n\r"), nil
}

// Install a package without confirmation.
func (z *Zypper) Install(ctx context.Context, pkgNames []string) error {
	args := []string{"install", "--no-confirm"}
	args = append(args, pkgNames...)

	cmd := z.command(ctx, false, true, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// command returns base zypper in quiet mode.
func (z *Zypper) command(ctx context.Context, quiet, sudo bool, args ...string) *exec.Cmd {
	if quiet {
		args = append([]string{"--quiet"}, args...)
	}

	name := "zypper"
	if sudo {
		name = "sudo"
		args = append([]string{"zypper"}, args...)
	}

	return exec.CommandContext(ctx, name, args...)
}
