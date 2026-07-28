package logs

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
)

// current is the active session, set by Open and cleared by Close.
var current *Session

// Session captures the stdout/stderr of a single process run to a timestamped log file.
type Session struct {
	file   *os.File
	logger *slog.Logger
	out    io.Writer
	err    io.Writer
}

// Logger returns the slog.Logger for explicit logging.
func (s *Session) Logger() *slog.Logger {
	return s.logger
}

// Stdout returns a writer that tees to the terminal stdout and the log file.
func Stdout() io.Writer {
	if current == nil {
		return os.Stdout
	}

	return current.out
}

// Stderr returns a writer that tees to the terminal stderr and the log file.
func Stderr() io.Writer {
	if current == nil {
		return os.Stderr
	}

	return current.err
}

// Open creates a timestamped log file under XDG_STATE_HOME/vlx/logs/.
// The command name is included in the filename for easier identification.
func Open(cmd string) (*Session, error) {
	dir := fsys.StatePath("vlx", "logs")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("logs: %w", err)
	}

	name := time.Now().Format("2006-01-02T15-04-05") + "_" + cmd + ".log"
	path := filepath.Join(dir, name)

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("logs: %w", err)
	}

	s := &Session{
		file:   file,
		logger: slog.New(slog.NewTextHandler(file, nil)),
		out:    io.MultiWriter(os.Stdout, file),
		err:    io.MultiWriter(os.Stderr, file),
	}

	_ = truncate(dir, 50)

	current = s
	return s, nil
}

// truncate keeps only the latest n log files in the directory,
// deleting older ones. Files are sorted by name (timestamp prefix).
func truncate(dir string, max int) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var logs []os.DirEntry
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".log" {
			logs = append(logs, e)
		}
	}

	if len(logs) <= max {
		return nil
	}

	sort.Slice(logs, func(i, j int) bool {
		return logs[i].Name() < logs[j].Name()
	})

	for _, e := range logs[:len(logs)-max] {
		_ = os.Remove(filepath.Join(dir, e.Name()))
	}

	return nil
}

// Close closes the log file and clears the active session.
func (s *Session) Close() error {
	current = nil
	return s.file.Close()
}
