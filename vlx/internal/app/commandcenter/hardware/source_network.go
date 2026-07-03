package hardware

import (
	"context"
	"fmt"

	"github.com/rlvelte/velinux/vlx/internal/core/printer"
)

func init() {
	register(&networkSource{})
}

type networkSource struct{}

func (s *networkSource) Name() string      { return "network" }
func (s *networkSource) Aliases() []string { return []string{"net", "interfaces"} }

func (s *networkSource) Run(ctx context.Context, p *printer.Printer, json bool) error {
	info, err := readNetwork()
	if err != nil {
		return fmt.Errorf("reading network info: %w", err)
	}

	if json {
		PrintJSON(info)
		return nil
	}

	if len(info.Interfaces) == 0 {
		p.Info("No network interfaces found")
		return nil
	}

	rows := make([][]string, 0, len(info.Interfaces))
	for _, iface := range info.Interfaces {
		status := "down"
		if iface.Up {
			status = "up"
		}
		rows = append(rows, []string{
			iface.Name,
			status,
			formatBytes(iface.RxBytes),
			formatBytes(iface.TxBytes),
		})
	}

	p.Table([]string{"Interface", "Status", "RX", "TX"}, rows)
	return nil
}
