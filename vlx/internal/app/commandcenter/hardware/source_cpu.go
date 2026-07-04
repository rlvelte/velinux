package hardware

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rlvelte/velinux/vlx/internal/core/printer"
)

func init() {
	register(&cpuSource{})
}

type cpuSource struct{}

func (s *cpuSource) Name() string      { return "cpu" }
func (s *cpuSource) Aliases() []string { return []string{"processor", "cpustats"} }

func (s *cpuSource) Run(ctx context.Context, p *printer.Printer, json bool) error {
	info, err := readCPU()
	if err != nil {
		return fmt.Errorf("reading CPU info: %w", err)
	}

	if json {
		PrintJSON(info)
		return nil
	}

	p.Table([]string{"Property", "Value"}, [][]string{
		{"Model", info.Model},
		{"Cores", strconv.Itoa(info.Cores)},
		{"Threads", strconv.Itoa(info.Threads)},
		{"Frequency", fmt.Sprintf("%.0f MHz", info.Freq)},
		{"Usage", fmt.Sprintf("%.1f%%", info.Usage)},
		{"Temperature", fmt.Sprintf("%.0f°C", info.Temp)},
	})
	return nil
}
