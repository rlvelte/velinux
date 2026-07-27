package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	"github.com/rlvelte/velinux/vlx/internal/core/guard"
)

type promptData struct {
	Prompt string `json:"prompt"`
}

func main() {
	prompt := "Password:"
	if len(os.Args) > 1 {
		prompt = os.Args[1]
	}

	if err := run(prompt); err != nil {
		fmt.Fprintf(os.Stderr, "vlxpass: %s\n", err)
		os.Exit(1)
	}
}

// run triggers a quickshell password popup and returns the entered password.
func run(prompt string) error {
	if err := guard.Binaries("quickshell"); err != nil {
		return err
	}

	base := dir()
	if err := os.MkdirAll(base, 0700); err != nil {
		return fmt.Errorf("su: creating runtime dir: %w", err)
	}

	clean(base)

	dir, err := os.MkdirTemp(base, "vlx-su-*")
	if err != nil {
		return fmt.Errorf("su: creating temp dir: %w", err)
	}
	defer os.RemoveAll(dir)

	data, err := json.Marshal(promptData{Prompt: prompt})
	if err != nil {
		return fmt.Errorf("su: marshalling prompt: %w", err)
	}

	if err := fsys.AtomicWrite(filepath.Join(dir, "prompt.json"), data); err != nil {
		return fmt.Errorf("su: writing prompt: %w", err)
	}

	if err := exec.Command("quickshell", "ipc", "call", "su", "vlxOpen", dir).Run(); err != nil {
		return fmt.Errorf("su: calling quickshell: %w", err)
	}

	resultPath := filepath.Join(dir, "result")
	if err := wait(resultPath, 120*time.Second); err != nil {
		return err
	}

	raw, err := os.ReadFile(resultPath)
	if err != nil {
		return fmt.Errorf("su: reading result: %w", err)
	}

	if len(raw) == 0 {
		return fmt.Errorf("su: cancelled")
	}
	
	fmt.Println(strings.TrimRight(string(raw), "\n"))
	return nil
}

func dir() string {
	if v := os.Getenv("XDG_RUNTIME_DIR"); v != "" {
		return filepath.Join(v, "vlx")
	}

	return filepath.Join("/dev/shm", fmt.Sprintf("vlx-%d", os.Getuid()))
}

func clean(base string) {
	entries, _ := os.ReadDir(base)
	cutoff := time.Now().Add(-10 * time.Minute)

	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), "vlx-su-") {
			continue
		}

		if info, err := e.Info(); err == nil && info.ModTime().Before(cutoff) {
			_ = os.RemoveAll(filepath.Join(base, e.Name()))
		}
	}
}

func wait(path string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return nil
		}

		time.Sleep(50 * time.Millisecond)
	}

	return fmt.Errorf("su: timed out waiting for password")
}
