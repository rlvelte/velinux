package hardware

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rlvelte/velinux/vlx/internal/core/printer"
)

func init() {
	register(&systemSource{})
}

type systemSource struct{}

func (s *systemSource) Name() string      { return "system" }
func (s *systemSource) Aliases() []string { return []string{"sys", "os", "kernel"} }

func (s *systemSource) Run(ctx context.Context, p *printer.Printer, json bool) error {
	info, err := readSystem()
	if err != nil {
		return fmt.Errorf("reading system info: %w", err)
	}

	if json {
		PrintJSON(info)
		return nil
	}

	rows := [][]string{
		{"Hostname", info.Hostname},
		{"OS", info.OS},
		{"Kernel", info.Kernel},
		{"Uptime", info.Uptime},
		{"Load Avg", fmt.Sprintf("%.2f %.2f %.2f", info.LoadAvg[0], info.LoadAvg[1], info.LoadAvg[2])},
		{"Processes", strconv.Itoa(info.Processes)},
	}

	p.Table([]string{"Property", "Value"}, rows)
	return nil
}
