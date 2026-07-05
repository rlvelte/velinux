package sources

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
)

// CPUInfo contains information about cpu.
type CPUInfo struct {
	Model   string  `json:"model"`
	Cores   int     `json:"cores"`
	Threads int     `json:"threads"`
	Freq    float64 `json:"freq_mhz"`
	MinFreq float64 `json:"min_freq_mhz,omitempty"`
	MaxFreq float64 `json:"max_freq_mhz,omitempty"`
	Usage   float64 `json:"usage_percent"`
	Temp    float64 `json:"temp_celsius,omitempty"`
}

// cpu information container.
type cpu struct{}

// Name of this source.
func (s *cpu) Name() string {
	return "cpu"
}

// Aliases that this source has.
func (s *cpu) Aliases() []string {
	return []string{"processor", "cpustats"}
}

// Run extracts all data for this source.
func (s *cpu) Run(p *printer.Printer) error {
	info, err := readCpu()
	if err != nil {
		return fmt.Errorf("reading CPU info: %w", err)
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

// readCpu aggregates all available data for cpu.
func readCpu() (*CPUInfo, error) {
	data, err := fsys.Proc("cpuinfo")
	if err != nil {
		return nil, err
	}

	info := &CPUInfo{}
	var modelNames []string
	coreIDs := make(map[string]bool)

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "model name") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				modelNames = append(modelNames, strings.TrimSpace(parts[1]))
			}
		}

		if strings.HasPrefix(line, "core id") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				coreIDs[strings.TrimSpace(parts[1])] = true
			}
		}

		if strings.HasPrefix(line, "cpu MHz") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				if f, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64); err == nil {
					info.Freq = f
				}
			}
		}
	}

	if len(modelNames) == 0 {
		modelNames = []string{"unknown"}
	}

	info.Model = modelNames[0]
	info.Threads = len(modelNames)
	info.Cores = len(coreIDs)
	if info.Cores == 0 {
		info.Cores = info.Threads
	}

	usage, err := readCpuUsage()
	if err == nil {
		info.Usage = usage
	}

	temp, err := readCpuTemp()
	if err == nil {
		info.Temp = temp
	}

	return info, nil
}

// readCpuUsage aggregates usage data for cpu.
func readCpuUsage() (float64, error) {
	data, err := fsys.Proc("stat")
	if err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	if !scanner.Scan() {
		return 0, fmt.Errorf("empty /proc/stat")
	}

	line := scanner.Text()
	if !strings.HasPrefix(line, "cpu ") {
		return 0, fmt.Errorf("unexpected /proc/stat format")
	}

	fields := strings.Fields(line)
	if len(fields) < 5 {
		return 0, fmt.Errorf("too few fields in /proc/stat")
	}

	var total uint64
	for i := 1; i < len(fields); i++ {
		val, _ := strconv.ParseUint(fields[i], 10, 64)
		total += val
	}

	idle, _ := strconv.ParseUint(fields[4], 10, 64)
	if len(fields) > 5 {
		iowait, _ := strconv.ParseUint(fields[5], 10, 64)
		idle += iowait
	}

	if total == 0 {
		return 0, nil
	}

	return float64(total-idle) / float64(total) * 100, nil
}

// readCpuTemp aggregates temperature data for cpu.
func readCpuTemp() (float64, error) {
	matches, err := filepath.Glob("/sys/class/thermal/thermal_zone*/temp")
	if err != nil || len(matches) == 0 {
		return 0, fmt.Errorf("no thermal zones")
	}

	for _, m := range matches {
		base := strings.TrimSuffix(m, "/temp")
		typeData, _ := os.ReadFile(filepath.Join(base, "type"))
		typ := strings.TrimSpace(string(typeData))

		if strings.Contains(strings.ToLower(typ), "cpu") || strings.Contains(strings.ToLower(typ), "core") || strings.Contains(strings.ToLower(typ), "x86") {
			val, err := fsys.ReadUInt64(m)
			if err != nil {
				continue
			}

			return float64(val) / 1000.0, nil
		}
	}

	val, err := fsys.ReadUInt64(matches[0])
	if err != nil {
		return 0, err
	}

	return float64(val) / 1000.0, nil
}
