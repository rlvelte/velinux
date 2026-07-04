package hardware

import (
	"context"
	"fmt"

	"github.com/rlvelte/velinux/vlx/internal/core/printer"
)

func init() {
	register(&tempSource{})
}

type tempSource struct{}

func (s *tempSource) Name() string      { return "temp" }
func (s *tempSource) Aliases() []string { return []string{"temperature", "sensors", "thermal"} }

func (s *tempSource) Run(ctx context.Context, p *printer.Printer, json bool) error {
	info, err := readTemps()
	if err != nil {
		return fmt.Errorf("reading temperature sensors: %w", err)
	}

	if json {
		PrintJSON(info)
		return nil
	}

	rows := make([][]string, 0, len(info.Sensors))
	for _, s := range info.Sensors {
		label := s.Label
		if label == "" {
			label = s.Name
		}
		rows = append(rows, []string{label, fmt.Sprintf("%.0f°C", s.Temp)})
	}

	p.Table([]string{"Sensor", "Temperature"}, rows)
	return nil
}
