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
		prompt = strings.Join(os.Args[1:], " ")
	}

	password, err := Show(prompt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "vlxpass: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(password)
}

// Show triggers a quickshell password popup and returns the entered password.
func Show(prompt string) (string, error) {
	if guard.Binaries("quickshell") != nil {
		return "", fmt.Errorf("pass: quickshell not found")
	}

	entries, err := os.ReadDir("/dev/shm")
	if err != nil {
		return "", fmt.Errorf("pass: reading /dev/shm: %w", err)
	}

	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "vlx-pass-") {
			os.RemoveAll(filepath.Join("/dev/shm", e.Name()))
		}
	}

	ts := fmt.Sprintf("%d", time.Now().UnixNano())
	dir := filepath.Join("/dev/shm", "vlx-pass-"+ts)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("pass: creating temp dir: %w", err)
	}
	defer os.RemoveAll(dir)

	data, err := json.Marshal(promptData{Prompt: prompt})
	if err != nil {
		return "", fmt.Errorf("pass: marshalling prompt: %w", err)
	}

	if err := fsys.AtomicWrite(filepath.Join(dir, "prompt.json"), data); err != nil {
		return "", fmt.Errorf("pass: writing prompt: %w", err)
	}

	if err := exec.Command("quickshell", "ipc", "call", "pass", "vlxOpen", dir).Run(); err != nil {
		return "", fmt.Errorf("pass: calling quickshell: %w", err)
	}

	resultPath := filepath.Join(dir, "result")
	if err := wait(resultPath, 120*time.Second); err != nil {
		return "", err
	}

	raw, err := os.ReadFile(resultPath)
	if err != nil {
		return "", fmt.Errorf("pass: reading result: %w", err)
	}

	if len(raw) == 0 {
		return "", fmt.Errorf("pass: cancelled")
	}

	return string(raw), nil
}

func wait(path string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return nil
		}
		<-ticker.C
	}

	return fmt.Errorf("pass: timed out waiting for password")
}
