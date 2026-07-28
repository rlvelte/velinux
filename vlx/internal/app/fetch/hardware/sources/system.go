package sources

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
)

// SystemInfo contains information about system.
type SystemInfo struct {
	Hostname  string     `json:"hostname"`
	OS        string     `json:"os"`
	Kernel    string     `json:"kernel"`
	Uptime    string     `json:"uptime"`
	UptimeSec float64    `json:"uptime_seconds"`
	LoadAvg   [3]float64 `json:"load_average"`
	Processes int        `json:"processes"`
}

// system information container.
type system struct{}

// Name of this source.
func (s *system) Name() string {
	return "system"
}

// Aliases that this source has.
func (s *system) Aliases() []string {
	return []string{"sys", "os", "kernel"}
}

// Run extracts all data for this source.
func (s *system) Run(p *printer.Printer) error {
	info, err := readSystem()
	if err != nil {
		return fmt.Errorf("reading system info: %w", err)
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

// readSystem aggregates all available data for system.
func readSystem() (*SystemInfo, error) {
	info := &SystemInfo{}

	hostname, err := os.Hostname()
	if err == nil {
		info.Hostname = hostname
	}

	if data, err := fsys.Proc("version"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) >= 3 {
			info.Kernel = fields[2]
		}
	}

	if data, err := fsys.Proc("uptime"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) > 0 {
			if sec, err := strconv.ParseFloat(fields[0], 64); err == nil {
				info.UptimeSec = sec
				info.Uptime = formatUptime(sec)
			}
		}
	}

	if data, err := fsys.Proc("loadavg"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) >= 5 {
			for i := 0; i < 3 && i < len(fields); i++ {
				info.LoadAvg[i], _ = strconv.ParseFloat(fields[i], 64)
			}
			if procsParts := strings.SplitN(fields[3], "/", 2); len(procsParts) == 2 {
				info.Processes, _ = strconv.Atoi(procsParts[1])
			}
		}
	}

	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				info.OS = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), `"`)
				break
			}
		}
	}

	return info, nil
}

func formatUptime(sec float64) string {
	days := int(sec / 86400)
	sec = float64(int(sec) % 86400)
	hours := int(sec / 3600)
	sec = float64(int(sec) % 3600)
	mins := int(sec / 60)

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, mins)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, mins)
	}
	return fmt.Sprintf("%dm", mins)
}
