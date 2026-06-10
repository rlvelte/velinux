package themes

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rvelte/vlx/internal/app/themes/templates"
	"gopkg.in/ini.v1"
)

func currentTheme() (Theme, error) {
	return loadTheme("current"), nil
}

func listThemes() ([]Theme, error) {
	entries, err := os.ReadDir(stateDir())

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Theme{}, nil
		}
		return nil, err
	}

	var themes []Theme
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "current" {
			continue
		}
		theme := loadTheme(entry.Name())
		themes = append(themes, theme)
	}

	return themes, nil
}

func loadTheme(from string) Theme {
	path := filepath.Join(stateDir(), from, "theme.conf")

	cfg, err := ini.Load(path)
	if err != nil {
		return Theme{Location: path, Name: path, Valid: false}
	}

	section := cfg.Section("theme")
	return Theme{
		Icon:     section.Key("icon").String(),
		Id:       section.Key("id").String(),
		Name:     section.Key("name").String(),
		Location: path,
		Valid:    true,
	}
}

func loadThemeData(name string) (templates.ThemeColorData, error) {
	path := filepath.Join(stateDir(), name, "theme.conf")

	cfg, err := ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
		AllowShadows:        false,
	}, path)

	if err != nil {
		return templates.ThemeColorData{}, fmt.Errorf("failed to load theme data for %q: %w", name, err)
	}

	s := cfg.Section("theme")

	d := templates.ThemeColorData{
		Primary:         s.Key("color_primary").String(),
		PrimaryDim:      s.Key("color_primary_dim").String(),
		PrimarySubtle:   s.Key("color_primary_subtle").String(),
		PrimaryMuted:    s.Key("color_primary_muted").String(),
		Secondary:       s.Key("color_secondary").String(),
		SecondaryDim:    s.Key("color_secondary_dim").String(),
		SecondaryLight:  s.Key("color_secondary_light").String(),
		Accent:          s.Key("color_accent").String(),
		AccentDim:       s.Key("color_accent_dim").String(),
		AccentLight:     s.Key("color_accent_light").String(),
		Base:            s.Key("color_base").String(),
		Mantle:          s.Key("color_mantle").String(),
		Crust:           s.Key("color_crust").String(),
		Surface0:        s.Key("color_surface0").String(),
		Surface1:        s.Key("color_surface1").String(),
		Surface2:        s.Key("color_surface2").String(),
		Text:            s.Key("color_text").String(),
		Subtext:         s.Key("color_subtext").String(),
		Overlay:         s.Key("color_overlay").String(),
		Muted:           s.Key("color_muted").String(),
		Success:         s.Key("color_success").String(),
		Warning:         s.Key("color_warning").String(),
		WarningSubtle:   s.Key("color_warning_subtle").String(),
		Error:           s.Key("color_error").String(),
		ErrorSubtle:     s.Key("color_error_subtle").String(),
		Info:            s.Key("color_info").String(),
		InfoSubtle:      s.Key("color_info_subtle").String(),
		OnPrimary:       s.Key("color_on_primary").String(),
		OnSecondary:     s.Key("color_on_secondary").String(),
		OnAccent:        s.Key("color_on_accent").String(),
		OnSurface:       s.Key("color_on_surface").String(),
		FontName:        s.Key("font_name").String(),
		FontNameHeading: s.Key("font_name_heading").String(),
		FontNameMono:    s.Key("font_name_mono").String(),
		FontSize:        s.Key("font_size").String(),
		FontSizeSmall:   s.Key("font_size_small").String(),
		FontSizeLarge:   s.Key("font_size_large").String(),
		FontSizeHeading: s.Key("font_size_heading").String(),
	}

	return d, nil
}
