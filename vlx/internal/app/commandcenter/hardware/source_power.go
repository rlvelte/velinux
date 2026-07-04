package hardware

import (
	"context"
	"fmt"

	"github.com/rlvelte/velinux/vlx/internal/core/printer"
)

func init() {
	register(&powerSource{})
}

type powerSource struct{}

func (s *powerSource) Name() string      { return "power" }
func (s *powerSource) Aliases() []string { return []string{"battery", "bat", "ac"} }

func (s *powerSource) Run(ctx context.Context, p *printer.Printer, json bool) error {
	info, err := readPower()
	if err != nil {
		return fmt.Errorf("reading power info: %w", err)
	}

	if json {
		PrintJSON(info)
		return nil
	}

	rows := make([][]string, 0, len(info.Supplies))
	for _, s := range info.Supplies {
		capacity := "N/A"
		if s.Capacity > 0 {
			capacity = fmt.Sprintf("%.0f%%", s.Capacity)
		}
		rows = append(rows, []string{s.Name, s.Type, s.Status, capacity})
	}

	p.Table([]string{"Supply", "Type", "Status", "Capacity"}, rows)
	return nil
}
