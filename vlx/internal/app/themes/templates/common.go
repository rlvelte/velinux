package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rvelte/vlx/internal/system/xdg"
)

type ThemeColorData struct {
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

// Generate creates the new theme files and writes them.
func Generate(data ThemeColorData) error {
	cfgDir := xdg.ConfigPath()
	generators := []struct {
		name string
		fn   func(ThemeColorData, string) error
	}{
		{"sway", GenerateSway},
		{"hypr", GenerateHypr},
		{"waybar", GenerateWaybar},
		{"rofi", GenerateRofi},
		{"kitty", GenerateKitty},
		{"mako", GenerateMako},
	}

	for _, g := range generators {
		if err := g.fn(data, cfgDir); err != nil {
			return fmt.Errorf("failed to generate %s config: %w", g.name, err)
		}
	}

	return nil
}

func toRgba(hex string) string {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) == 6 {
		hex += "ff"
	}
	return hex
}

func write(path, tmpl string, data ThemeColorData) error {
	t, err := template.New("").Funcs(template.FuncMap{
		"toRgba": toRgba,
	}).Parse(tmpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, buf.Bytes(), 0o644)
}
