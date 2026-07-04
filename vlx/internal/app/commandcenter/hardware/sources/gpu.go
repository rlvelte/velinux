package sources

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rlvelte/velinux/vlx/internal/core/format"
	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
)

// GPUInfo contains information about gpu.
type GPUInfo struct {
	Vendor string  `json:"vendor"`
	Model  string  `json:"model"`
	Driver string  `json:"driver,omitempty"`
	Memory uint64  `json:"memory_bytes,omitempty"`
	Temp   float64 `json:"temp_celsius,omitempty"`
	Usage  float64 `json:"usage_percent,omitempty"`
}

// gpu information container.
type gpu struct{}

// Name of this source.
func (s *gpu) Name() string { return "gpu" }

// Aliases that this source has.
func (s *gpu) Aliases() []string { return []string{"graphics", "video"} }

// Run extracts all data for this source.
func (s *gpu) Run(p *printer.Printer) error {
	info, err := readGpu()
	if err != nil {
		return fmt.Errorf("reading GPU info: %w", err)
	}

	rows := [][]string{
		{"Vendor", info.Vendor},
		{"Model", info.Model},
	}
	if info.Memory > 0 {
		rows = append(rows, []string{"Memory", format.Bytes(info.Memory)})
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

// readGpu aggregates all available data for gpu.
func readGpu() (*GPUInfo, error) {
	matches, err := filepath.Glob("/sys/class/drm/card*/device/vendor")
	if err != nil {
		return nil, err
	}

	type gpuCard struct {
		dir    string
		vendor string
		vram   uint64
	}

	var cards []gpuCard
	for _, vPath := range matches {
		data, err := os.ReadFile(vPath)
		if err != nil {
			continue
		}
		vendorID := strings.TrimSpace(string(data))

		vendor := readGpuVendor(vendorID)
		if vendor == "" {
			continue
		}

		devDir := filepath.Dir(vPath)
		vram, _ := fsys.ReadUInt64(filepath.Join(devDir, "mem_info_vram_total"))
		cards = append(cards, gpuCard{dir: devDir, vendor: vendor, vram: vram})
	}

	if len(cards) == 0 {
		return nil, fmt.Errorf("no GPU found")
	}

	best := cards[0]
	for _, c := range cards[1:] {
		if c.vram > best.vram {
			best = c
		}
	}

	info := &GPUInfo{Vendor: best.vendor}

	if name, err := os.ReadFile(filepath.Join(best.dir, "product_name")); err == nil {
		info.Model = strings.TrimSpace(string(name))
	} else {
		if id, err := os.ReadFile(filepath.Join(best.dir, "device")); err == nil {
			info.Model = best.vendor + " GPU (" + strings.TrimSpace(string(id)) + ")"
		} else {
			info.Model = best.vendor + " GPU"
		}
	}

	if usage, err := fsys.ReadInt64(filepath.Join(best.dir, "gpu_busy_percent")); err == nil {
		info.Usage = float64(usage)
	}

	if total, err := fsys.ReadUInt64(filepath.Join(best.dir, "mem_info_vram_total")); err == nil {
		if used, err := fsys.ReadUInt64(filepath.Join(best.dir, "mem_info_vram_used")); err == nil {
			info.Memory = used
			_ = total
		}
	}

	info.Temp = readGpuTemp()
	return info, nil
}

// readGpuTemp aggregates temperature data for gpu.
func readGpuTemp() float64 {
	matches, err := filepath.Glob("/sys/class/hwmon/hwmon*/name")
	if err != nil {
		return 0
	}

	for _, namePath := range matches {
		data, _ := os.ReadFile(namePath)
		name := strings.TrimSpace(string(data))
		if name != "amdgpu" && !strings.HasPrefix(name, "nvidia") {
			continue
		}

		hwmonDir := filepath.Dir(namePath)
		val, err := fsys.ReadInt64(filepath.Join(hwmonDir, "temp1_input"))
		if err != nil {
			return 0
		}
		return float64(val) / 1000.0
	}

	return 0
}

// readGpuVendor resolves vendor name for gpu.
func readGpuVendor(id string) string {
	switch strings.TrimSpace(id) {
	case "0x10de", "10de":
		return "NVIDIA"
	case "0x1002", "1002":
		return "AMD"
	case "0x8086", "8086":
		return "Intel"
	case "0x1ae0", "1ae0":
		return "Google"
	default:
		return ""
	}
}
