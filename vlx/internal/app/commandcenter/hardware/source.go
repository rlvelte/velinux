package hardware

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rlvelte/velinux/vlx/internal/core/printer"
)

type HardwareSource interface {
	Name() string
	Aliases() []string
	Run(ctx context.Context, p *printer.Printer, json bool) error
}

var sources []HardwareSource

func register(s HardwareSource) {
	sources = append(sources, s)
}

func Find(name string) HardwareSource {
	for _, s := range sources {
		if s.Name() == name {
			return s
		}
		for _, a := range s.Aliases() {
			if a == name {
				return s
			}
		}
	}
	return nil
}

func ListSources(p *printer.Printer) {
	rows := make([][]string, 0, len(sources))
	for _, s := range sources {
		aliases := ""
		if len(s.Aliases()) > 0 {
			aliases = strings.Join(s.Aliases(), ", ")
		}
		rows = append(rows, []string{s.Name(), aliases})
	}
	p.Table([]string{"Source", "Aliases"}, rows)
}

func PrintJSON(v any) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf(`{"error":"%s"}`, err.Error())
		fmt.Println()
		return
	}
	fmt.Println(string(data))
}
