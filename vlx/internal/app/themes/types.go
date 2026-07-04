package themes

import (
	"encoding/json"
	"fmt"
)

type Theme struct {
	Icon      string `json:"icon"`
	Logo      string `json:"logo"`
	Id        string `json:"id"`
	Name      string `json:"name"`
	Wallpaper string `json:"wallpaper"`
	Active    bool   `json:"active"`
	Path      string
}

func decodeTheme(_ string, path string, data []byte) (*Theme, error) {
	var t Theme
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}

	t.Path = path
	if t.Id == "" {
		return nil, fmt.Errorf("theme %s has no id", path)
	}

	return &t, nil
}

type ThemeContent struct {
	Primary         string `json:"color_primary"`
	PrimaryDim      string `json:"color_primary_dim"`
	PrimarySubtle   string `json:"color_primary_subtle"`
	PrimaryMuted    string `json:"color_primary_muted"`
	Secondary       string `json:"color_secondary"`
	SecondaryDim    string `json:"color_secondary_dim"`
	SecondaryLight  string `json:"color_secondary_light"`
	Accent          string `json:"color_accent"`
	AccentDim       string `json:"color_accent_dim"`
	AccentLight     string `json:"color_accent_light"`
	Base            string `json:"color_base"`
	Mantle          string `json:"color_mantle"`
	Crust           string `json:"color_crust"`
	Surface0        string `json:"color_surface0"`
	Surface1        string `json:"color_surface1"`
	Surface2        string `json:"color_surface2"`
	Text            string `json:"color_text"`
	Subtext         string `json:"color_subtext"`
	Overlay         string `json:"color_overlay"`
	Muted           string `json:"color_muted"`
	Success         string `json:"color_success"`
	Warning         string `json:"color_warning"`
	WarningSubtle   string `json:"color_warning_subtle"`
	Error           string `json:"color_error"`
	ErrorSubtle     string `json:"color_error_subtle"`
	Info            string `json:"color_info"`
	InfoSubtle      string `json:"color_info_subtle"`
	OnPrimary       string `json:"color_on_primary"`
	OnSecondary     string `json:"color_on_secondary"`
	OnAccent        string `json:"color_on_accent"`
	OnSurface       string `json:"color_on_surface"`
	FontName        string `json:"font_name"`
	FontNameHeading string `json:"font_name_heading"`
	FontNameMono    string `json:"font_name_mono"`
	FontSize        string `json:"font_size"`
	FontSizeSmall   string `json:"font_size_small"`
	FontSizeLarge   string `json:"font_size_large"`
	FontSizeHeading string `json:"font_size_heading"`
}

func decodeThemeContent(_ string, _ string, data []byte) (*ThemeContent, error) {
	var tc ThemeContent
	if err := json.Unmarshal(data, &tc); err != nil {
		return nil, err
	}

	return &tc, nil
}
