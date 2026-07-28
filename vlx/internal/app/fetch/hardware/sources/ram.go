package sources

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
)

// RamInfo contains information about ram.
type RamInfo struct {
	Total     uint64 `json:"total_bytes"`
	Used      uint64 `json:"used_bytes"`
	Free      uint64 `json:"free_bytes"`
	Available uint64 `json:"available_bytes"`
	SwapTotal uint64 `json:"swap_total_bytes"`
	SwapUsed  uint64 `json:"swap_used_bytes"`
	SwapFree  uint64 `json:"swap_free_bytes"`
}

// ram information container.
type ram struct{}

// Name of this source.
func (s *ram) Name() string {
	return "memory"
}

// Aliases that this source has.
func (s *ram) Aliases() []string {
	return []string{"mem", "ram"}
}

// Run extracts all data for this source.
func (s *ram) Run(p *printer.Printer) error {
	info, err := readRam()
	if err != nil {
		return fmt.Errorf("reading memory info: %w", err)
	}

	p.Table([]string{"Property", "Value"}, [][]string{
		{"Total", FormatBytes(info.Total)},
		{"Used", FormatBytes(info.Used)},
		{"Free", FormatBytes(info.Free)},
		{"Available", FormatBytes(info.Available)},
		{"Swap Total", FormatBytes(info.SwapTotal)},
		{"Swap Used", FormatBytes(info.SwapUsed)},
		{"Swap Free", FormatBytes(info.SwapFree)},
	})
	return nil
}

// readRam aggregates all available data for ram.
func readRam() (*RamInfo, error) {
	data, err := fsys.Proc("meminfo")
	if err != nil {
		return nil, err
	}

	info := &RamInfo{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		fields := strings.Fields(strings.TrimSpace(parts[1]))
		if len(fields) == 0 {
			continue
		}

		val, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			continue
		}

		if len(fields) > 1 && fields[1] == "kB" {
			val *= 1024
		}

		switch key {
		case "MemTotal":
			info.Total = val
		case "MemFree":
			info.Free = val
		case "MemAvailable":
			info.Available = val
		case "SwapTotal":
			info.SwapTotal = val
		case "SwapFree":
			info.SwapFree = val
		}
	}

	info.Used = info.Total - info.Available
	info.SwapUsed = info.SwapTotal - info.SwapFree

	return info, nil
}
