package themes

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rvelte/vlx/internal/app/themes/templates"
	"github.com/rvelte/vlx/internal/core/cmd"
)

func switchTheme(query string) error {
	theme := loadTheme(query)
	if !theme.Valid {
		themes, err := listThemes()
		if err != nil {
			return fmt.Errorf("theme %q not found or invalid", query)
		}

		for _, t := range themes {
			if t.Name == query || t.Id == query {
				theme = t
				break
			}
		}

		if !theme.Valid {
			return fmt.Errorf("theme %q not found or invalid", query)
		}
	}

	themeDir := filepath.Base(filepath.Dir(theme.Location))

	if err := copyThemeToCurrent(themeDir); err != nil {
		return fmt.Errorf("failed to update current theme: %w", err)
	}

	data, err := loadThemeData(themeDir)
	if err != nil {
		return err
	}

	if err := templates.Generate(data); err != nil {
		return err
	}

	cmd.New("swaymsg", "reload").RunParallel()
	cmd.New("hyprctl", "reload").RunParallel()
	cmd.New("makoctl", "reload").RunParallel()

	cmd.New("pkill", "-USR1", "kitty").RunParallel()
	cmd.New("pkill", "-RTMIN+1", "waybar").RunParallel()

	return nil
}

func copyThemeToCurrent(name string) error {
	srcDir := filepath.Join(stateDir(), name)
	dstDir := filepath.Join(stateDir(), "current")

	if err := os.RemoveAll(dstDir); err != nil {
		return err
	}
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return err
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		src := filepath.Join(srcDir, entry.Name())
		dst := filepath.Join(dstDir, entry.Name())
		if err := copyFile(src, dst); err != nil {
			return err
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Close()
}
