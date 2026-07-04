package launcher

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// Entry represents a parsed desktop entry.
type Entry struct {
	ID        string // File basename without .desktop
	Name      string
	Comment   string
	Icon      string
	Exec      string
	NoDisplay bool
	Hidden    bool
	Terminal  bool
}

// Default scan directories in priority order.
var scanDirs = func() []string {
	home, _ := os.UserHomeDir()
	return []string{
		filepath.Join(home, ".local/share/applications"),
		"/usr/share/applications",
		"/usr/local/share/applications",
		filepath.Join(home, ".local/share/flatpak/exports/share/applications"),
		"/var/lib/flatpak/exports/share/applications",
		filepath.Join(home, ".config/autostart"),
	}
}()

// Scan walks the desktop entry directories and returns a slice of entries.
func Scan() []Entry {
	var entries []Entry
	seen := make(map[string]bool)

	for _, dir := range scanDirs {
		_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(path, ".desktop") {
				return nil
			}
			id := strings.TrimSuffix(filepath.Base(path), ".desktop")
			if seen[id] {
				return nil
			}

			ent := parse(path)
			if ent == nil {
				return nil
			}
			ent.ID = id
			seen[id] = true
			entries = append(entries, *ent)
			return nil
		})
	}

	return entries
}

// parse does a fast line-by-line read of a .desktop file.
// It stops scanning a file once all required fields are found.
func parse(path string) *Entry {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	ent := &Entry{}
	var inDesktopEntry bool
	var got int
	const need = 6

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if line == "[Desktop Entry]" {
			inDesktopEntry = true
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			if got >= need {
				break
			}
			inDesktopEntry = false
			continue
		}

		if !inDesktopEntry {
			continue
		}

		if strings.HasPrefix(line, "Name=") {
			ent.Name = strings.TrimPrefix(line, "Name=")
			got++
		} else if strings.HasPrefix(line, "Comment=") {
			ent.Comment = strings.TrimPrefix(line, "Comment=")
			got++
		} else if strings.HasPrefix(line, "Icon=") {
			ent.Icon = strings.TrimPrefix(line, "Icon=")
			got++
		} else if strings.HasPrefix(line, "Exec=") {
			ent.Exec = strings.TrimPrefix(line, "Exec=")
			got++
		} else if strings.HasPrefix(line, "NoDisplay=true") {
			ent.NoDisplay = true
			got++
		} else if strings.HasPrefix(line, "Hidden=true") {
			ent.Hidden = true
			got++
		} else if strings.HasPrefix(line, "Terminal=true") {
			ent.Terminal = true
		}

		if got >= need {
			break
		}
	}

	if ent.Name == "" || ent.Exec == "" {
		return nil
	}

	return ent
}
