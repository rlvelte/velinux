package themes

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

type target struct {
	template string
	output   func(cfgRoot string) string
}

var targets = []target{
	{"sway", func(d string) string { return filepath.Join(d, "sway", "config.d", "40-theme.conf") }},
	{"hypr", func(d string) string { return filepath.Join(d, "hypr", "config.d", "40-theme.conf") }},
	{"mangowm", func(d string) string { return filepath.Join(d, "mango", "config.d", "40-theme.conf") }},
	{"waybar", func(d string) string { return filepath.Join(d, "waybar", "theme.css") }},
	{"rofi", func(d string) string { return filepath.Join(d, "rofi", "theme.rasi") }},
	{"kitty", func(d string) string { return filepath.Join(d, "kitty", "theme.conf") }},
	{"mako", func(d string) string { return filepath.Join(d, "mako", "theme.conf") }},
}

// GenerateAll takes all target definitions and generate a new theme file.
func GenerateAll(from ThemeContent) error {
	for _, g := range targets {
		targetPath := g.output(fsys.ConfigPath())

		if err := render(targetPath, g.template, from); err != nil {
			return fmt.Errorf("failed to generate %s config: %w", g.template, err)
		}
	}

	return nil
}

func render(target, tmplName string, content ThemeContent) error {
	tmplStr, err := load(tmplName)
	if err != nil {
		return err
	}

	tmpl, err := template.New("").Funcs(template.FuncMap{
		"toRgba": rgba,
	}).Parse(tmplStr)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, content); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}

	return os.WriteFile(target, buf.Bytes(), 0644)
}

func load(name string) (string, error) {
	data, err := templateFS.ReadFile("templates/" + name + ".tmpl")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func rgba(hex string) string {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) == 6 {
		hex += "ff"
	}
	return hex
}
