package launcher

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
)

var iconSizes = []string{"scalable", "48x48", "64x64", "32x32", "128x128", "256x256", "24x24", "16x16", "symbolic"}
var iconExts = []string{".png", ".svg", ".xpm"}

var iconDirs = func() []string {
	dirs := []string{fsys.DataPath("icons")}
	for _, d := range fsys.DataDirs() {
		dirs = append(dirs, filepath.Join(d, "icons"))
	}

	return dirs
}()

var appDirs = func() []string {
	home, _ := os.UserHomeDir()

	dirs := []string{fsys.DataPath("applications")}
	for _, d := range fsys.DataDirs() {
		dirs = append(dirs, filepath.Join(d, "applications"))
	}

	dirs = append(dirs,
		filepath.Join(home, ".local/share/flatpak/exports/share/applications"),
		"/var/lib/flatpak/exports/share/applications",
		filepath.Join(home, ".config/autostart"),
	)

	return dirs
}()

// Entry represents a parsed desktop entry.
type Entry struct {
	ID        string
	Name      string
	Comment   string
	Icon      string
	Exec      string
	NoDisplay bool
	Hidden    bool
}

// Scan walks the desktop entry directories and returns a slice of entries.
func Scan() []Entry {
	var entries []Entry
	seen := make(map[string]bool)

	for _, dir := range appDirs {
		_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(path, ".desktop") {
				return nil
			}

			id := strings.TrimSuffix(filepath.Base(path), ".desktop")
			if seen[id] {
				return nil
			}

			ent := parseApp(path)
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

// parseApp does a fast line-by-line read of a .desktop file.
func parseApp(path string) *Entry {
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

// parseIcon resolves an XDG icon name to an absolute path.
func parseIcon(icon string) string {
	if icon == "" || filepath.IsAbs(icon) {
		return icon
	}

	for _, root := range iconDirs {
		entries, err := os.ReadDir(root)
		if err != nil {
			continue
		}

		for _, e := range entries {
			if !e.IsDir() {
				continue
			}

			themeDir := filepath.Join(root, e.Name())
			for _, size := range iconSizes {
				for _, ext := range iconExts {
					path := filepath.Join(themeDir, size, "apps", icon+ext)
					if _, err := os.Stat(path); err == nil {
						return path
					}
				}
			}
		}
	}

	for _, ext := range iconExts {
		path := filepath.Join("/usr/share/pixmaps", icon+ext)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return icon
}
