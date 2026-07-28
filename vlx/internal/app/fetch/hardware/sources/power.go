package sources

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
)

// PowerInfo contains information about power.
type PowerInfo struct {
	Supplies []PowerSupply `json:"supplies"`
}

type PowerSupply struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Status   string  `json:"status"`
	Capacity float64 `json:"capacity_percent,omitempty"`
}

// power information container.
type power struct{}

// Name of this source.
func (s *power) Name() string {
	return "power"
}

// Aliases that this source has.
func (s *power) Aliases() []string {
	return []string{"battery", "bat", "ac"}
}

// Run extracts all data for this source.
func (s *power) Run(p *printer.Printer) error {
	info, err := readPower()
	if err != nil {
		return fmt.Errorf("reading power info: %w", err)
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

// readPower aggregates all available data for power.
func readPower() (*PowerInfo, error) {
	matches, err := filepath.Glob("/sys/class/power_supply/*")
	if err != nil {
		return nil, err
	}

	info := &PowerInfo{}
	for _, dir := range matches {
		name := filepath.Base(dir)

		typeData, _ := os.ReadFile(filepath.Join(dir, "type"))
		typ := strings.TrimSpace(string(typeData))

		statusData, _ := os.ReadFile(filepath.Join(dir, "status"))
		status := strings.TrimSpace(string(statusData))

		var capacity float64
		if capVal, err := fsys.ReadInt64(filepath.Join(dir, "capacity")); err == nil {
			capacity = float64(capVal)
		}

		info.Supplies = append(info.Supplies, PowerSupply{
			Name:     name,
			Type:     typ,
			Status:   status,
			Capacity: capacity,
		})
	}

	if len(info.Supplies) == 0 {
		return nil, fmt.Errorf("no power supplies found")
	}

	return info, nil
}
