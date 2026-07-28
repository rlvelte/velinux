package sources

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
)

// TempInfo contains information about temperature.
type TempInfo struct {
	Sensors []TempSensor `json:"sensors"`
}

type TempSensor struct {
	Name  string  `json:"name"`
	Label string  `json:"label,omitempty"`
	Temp  float64 `json:"temp_celsius"`
}

// temperature information container.
type temperature struct{}

// Name of this source.
func (s *temperature) Name() string {
	return "temp"
}

// Aliases that this source has.
func (s *temperature) Aliases() []string {
	return []string{"temperature", "sensors", "thermal"}
}

// Run extracts all data for this source.
func (s *temperature) Run(p *printer.Printer) error {
	info, err := readTemperature()
	if err != nil {
		return fmt.Errorf("reading temperature sensors: %w", err)
	}

	rows := make([][]string, 0, len(info.Sensors))
	for _, s := range info.Sensors {
		label := s.Label
		if label == "" {
			label = s.Name
		}
		rows = append(rows, []string{label, fmt.Sprintf("%.0f°C", s.Temp)})
	}

	p.Table([]string{"Sensor", "Temperature"}, rows)
	return nil
}

// readTemperature aggregates all available data for temperature.
func readTemperature() (*TempInfo, error) {
	info := &TempInfo{}

	sensors, err := readTemperatureHwmon()
	if err == nil {
		info.Sensors = append(info.Sensors, sensors...)
	}

	thermal, err := readTemperatureZone()
	if err == nil {
		info.Sensors = append(info.Sensors, thermal...)
	}

	if len(info.Sensors) == 0 {
		return nil, fmt.Errorf("no temperature sensors found")
	}

	return info, nil
}

// readTemperatureHwmon aggregates hwmon sensor data for temperature.
func readTemperatureHwmon() ([]TempSensor, error) {
	matches, err := filepath.Glob("/sys/class/hwmon/hwmon*")
	if err != nil {
		return nil, err
	}

	var sensors []TempSensor
	for _, dir := range matches {
		nameData, err := os.ReadFile(filepath.Join(dir, "name"))
		if err != nil {
			continue
		}
		name := strings.TrimSpace(string(nameData))

		inputs, _ := filepath.Glob(filepath.Join(dir, "temp*_input"))
		for _, inputPath := range inputs {
			base := strings.TrimSuffix(inputPath, "_input")
			val, err := fsys.ReadInt64(inputPath)
			if err != nil {
				continue
			}

			label := ""
			if labelData, err := os.ReadFile(base + "_label"); err == nil {
				label = strings.TrimSpace(string(labelData))
			}

			sensors = append(sensors, TempSensor{
				Name:  name,
				Label: label,
				Temp:  float64(val) / 1000.0,
			})
		}
	}

	return sensors, nil
}

// readTemperatureZone aggregates thermal zone data for temperature.
func readTemperatureZone() ([]TempSensor, error) {
	matches, err := filepath.Glob("/sys/class/thermal/thermal_zone*")
	if err != nil {
		return nil, err
	}

	var zones []TempSensor
	for _, dir := range matches {
		val, err := fsys.ReadInt64(filepath.Join(dir, "temp"))
		if err != nil {
			continue
		}

		typeData, _ := os.ReadFile(filepath.Join(dir, "type"))
		typ := strings.TrimSpace(string(typeData))

		zones = append(zones, TempSensor{
			Name: typ,
			Temp: float64(val) / 1000.0,
		})
	}

	return zones, nil
}
