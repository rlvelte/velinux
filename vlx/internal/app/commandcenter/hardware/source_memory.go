package hardware

import (
	"context"
	"fmt"

	"github.com/rlvelte/velinux/vlx/internal/core/printer"
)

func init() {
	register(&memorySource{})
}

type memorySource struct{}

func (s *memorySource) Name() string      { return "memory" }
func (s *memorySource) Aliases() []string { return []string{"mem", "ram"} }

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func (s *memorySource) Run(ctx context.Context, p *printer.Printer, json bool) error {
	info, err := readMemory()
	if err != nil {
		return fmt.Errorf("reading memory info: %w", err)
	}

	if json {
		PrintJSON(info)
		return nil
	}

	p.Table([]string{"Property", "Value"}, [][]string{
		{"Total", formatBytes(info.Total)},
		{"Used", formatBytes(info.Used)},
		{"Free", formatBytes(info.Free)},
		{"Available", formatBytes(info.Available)},
		{"Swap Total", formatBytes(info.SwapTotal)},
		{"Swap Used", formatBytes(info.SwapUsed)},
		{"Swap Free", formatBytes(info.SwapFree)},
	})
	return nil
}
