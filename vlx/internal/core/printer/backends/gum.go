package backends

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Gum is a backend that uses the gum command-line tool.
type Gum struct{}

// Available reports whether gum is installed.
func (g *Gum) Available() bool {
	_, err := exec.LookPath("gum")
	return err == nil
}

// Info prints an info message.
func (g *Gum) Info(msg string) {
	g.log("info", msg)
}

// Success prints a success message.
func (g *Gum) Success(msg string) {
	g.log("info", "✓ "+msg)
}

// Warn prints a warning message.
func (g *Gum) Warn(msg string) {
	g.log("warn", msg)
}

// Error prints an error message.
func (g *Gum) Error(msg string) {
	g.log("error", msg)
}

// Header signals a new section.
func (g *Gum) Header(msg string) {
	cmd := exec.Command("gum", "style", "--bold", msg)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Println("==>", msg)
	}
}

// Table prints data in a tabular format.
func (g *Gum) Table(headers []string, rows [][]string) {
	args := []string{"table", "--print"}
	if len(headers) > 0 {
		args = append(args, "--columns", joinCSV(headers))
	}

	cmd := exec.Command("gum", args...)
	cmd.Stdin = strings.NewReader(rowsToCSV(rows))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		(&Basic{}).Table(headers, rows)
	}
}

// Confirm shows a simple confirmation dialog.
func (g *Gum) Confirm(msg string) bool {
	cmd := exec.Command("gum", "confirm", msg)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run() == nil
}

// Spinner shows a progress indicator.
func (g *Gum) Spinner(label string, fn func() error) error {
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	done := make(chan error, 1)

	go func() {
		done <- fn()
	}()

	i := 0
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case err := <-done:
			fmt.Fprintf(os.Stderr, "\r%s %s", label, "done")
			if err != nil {
				fmt.Fprintf(os.Stderr, "\r%s %s\n", label, "failed")
			} else {
				fmt.Fprintf(os.Stderr, "\r%s %s\n", label, "done")
			}
			return err
		case <-ticker.C:
			fmt.Fprintf(os.Stderr, "\r%s %s", spinner[i%len(spinner)], label)
			i++
		}
	}
}

func (g *Gum) log(level, msg string) {
	cmd := exec.Command("gum", "log", "--level", level, msg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(msg)
	}
}

func joinCSV(fields []string) string {
	quoted := make([]string, len(fields))
	for i, f := range fields {
		if strings.ContainsAny(f, ",\"\n") {
			quoted[i] = `"` + strings.ReplaceAll(f, `"`, `""`) + `"`
		} else {
			quoted[i] = f
		}
	}
	return strings.Join(quoted, ",")
}

func rowsToCSV(rows [][]string) string {
	lines := make([]string, len(rows))
	for i, row := range rows {
		lines[i] = joinCSV(row)
	}
	return strings.Join(lines, "\n")
}
