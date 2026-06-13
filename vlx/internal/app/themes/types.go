package themes

import (
	"fmt"

	"gopkg.in/ini.v1"
)

type Theme struct {
	Icon      string
	Id        string
	Name      string
	Wallpaper string
	Path      string
}

func decodeTheme(_ string, path string, data []byte) (*Theme, error) {
	cfg, err := ini.LoadSources(ini.LoadOptions{SpaceBeforeInlineComment: true}, data)
	if err != nil {
		return nil, err
	}

	s := cfg.Section("theme")
	t := &Theme{
		Icon:      s.Key("icon").String(),
		Id:        s.Key("id").String(),
		Name:      s.Key("name").String(),
		Wallpaper: s.Key("wallpaper").String(),
		Path:      path,
	}
	if t.Id == "" {
		return nil, fmt.Errorf("theme %s has no id", path)
	}
	return t, nil
}

type ThemeContent struct {
	Primary         string
	PrimaryDim      string
	PrimarySubtle   string
	PrimaryMuted    string
	Secondary       string
	SecondaryDim    string
	SecondaryLight  string
	Accent          string
	AccentDim       string
	AccentLight     string
	Base            string
	Mantle          string
	Crust           string
	Surface0        string
	Surface1        string
	Surface2        string
	Text            string
	Subtext         string
	Overlay         string
	Muted           string
	Success         string
	Warning         string
	WarningSubtle   string
	Error           string
	ErrorSubtle     string
	Info            string
	InfoSubtle      string
	OnPrimary       string
	OnSecondary     string
	OnAccent        string
	OnSurface       string
	FontName        string
	FontNameHeading string
	FontNameMono    string
	FontSize        string
	FontSizeSmall   string
	FontSizeLarge   string
	FontSizeHeading string
}

func decodeThemeContent(_ string, _ string, data []byte) (*ThemeContent, error) {
	cfg, err := ini.LoadSources(ini.LoadOptions{SpaceBeforeInlineComment: true}, data)
	if err != nil {
		return nil, err
	}
	s := cfg.Section("theme")
	return &ThemeContent{
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
	}, nil
}
