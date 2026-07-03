package hardware

import (
	"context"
	"fmt"

	"github.com/rlvelte/velinux/vlx/internal/core/printer"
)

func init() {
	register(&gpuSource{})
}

type gpuSource struct{}

func (s *gpuSource) Name() string      { return "gpu" }
func (s *gpuSource) Aliases() []string { return []string{"graphics", "video"} }

func (s *gpuSource) Run(ctx context.Context, p *printer.Printer, json bool) error {
	info, err := readGPU()
	if err != nil {
		return fmt.Errorf("reading GPU info: %w", err)
	}

	if json {
		PrintJSON(info)
		return nil
	}

	rows := [][]string{
		{"Vendor", info.Vendor},
		{"Model", info.Model},
	}
	if info.Memory > 0 {
		rows = append(rows, []string{"Memory", formatBytes(info.Memory)})
	}
	if info.Temp > 0 {
		rows = append(rows, []string{"Temperature", fmt.Sprintf("%.0f°C", info.Temp)})
	}
	if info.Usage > 0 {
		rows = append(rows, []string{"Usage", fmt.Sprintf("%.0f%%", info.Usage)})
	}

	p.Table([]string{"Property", "Value"}, rows)
	return nil
}
