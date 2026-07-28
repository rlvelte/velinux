package sources

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
)

// NetworkInfo contains information about network.
type NetworkInfo struct {
	Interfaces []NetInterface `json:"interfaces"`
}

type NetInterface struct {
	Name    string `json:"name"`
	Up      bool   `json:"up"`
	RxBytes uint64 `json:"rx_bytes"`
	TxBytes uint64 `json:"tx_bytes"`
}

// network information container.
type network struct{}

// Name of this source.
func (s *network) Name() string {
	return "network"
}

// Aliases that this source has.
func (s *network) Aliases() []string {
	return []string{"net", "interfaces"}
}

// Run extracts all data for this source.
func (s *network) Run(p *printer.Printer) error {
	info, err := readNetwork()
	if err != nil {
		return fmt.Errorf("reading network info: %w", err)
	}

	if len(info.Interfaces) == 0 {
		p.Success("No network interfaces found")
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
			FormatBytes(iface.RxBytes),
			FormatBytes(iface.TxBytes),
		})
	}

	p.Table([]string{"Interface", "Status", "RX", "TX"}, rows)
	return nil
}

// readNetwork aggregates all available data for network.
func readNetwork() (*NetworkInfo, error) {
	info := &NetworkInfo{}

	interfaces, err := filepath.Glob("/sys/class/net/*")
	if err != nil {
		return nil, err
	}

	for _, ifPath := range interfaces {
		name := filepath.Base(ifPath)

		if strings.HasPrefix(name, "lo") {
			continue
		}

		operData, err := os.ReadFile(filepath.Join(ifPath, "operstate"))
		up := err == nil && strings.TrimSpace(string(operData)) == "up"

		var rxBytes, txBytes uint64
		if up {
			rxBytes, _ = fsys.ReadUInt64(filepath.Join(ifPath, "statistics", "rx_bytes"))
			txBytes, _ = fsys.ReadUInt64(filepath.Join(ifPath, "statistics", "tx_bytes"))
		}

		info.Interfaces = append(info.Interfaces, NetInterface{
			Name:    name,
			Up:      up,
			RxBytes: rxBytes,
			TxBytes: txBytes,
		})
	}

	return info, nil
}
