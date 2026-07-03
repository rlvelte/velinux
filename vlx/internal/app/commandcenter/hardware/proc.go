package hardware

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func readProc(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join("/proc", name))
}

func readSys(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join("/sys", name))
}

func readInt(path string) (int64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
}

func readUint64(path string) (uint64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
}

func readCPU() (*CPUInfo, error) {
	data, err := readProc("cpuinfo")
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

	usage, err := cpuUsage()
	if err == nil {
		info.Usage = usage
	}

	temp, err := cpuTemp()
	if err == nil {
		info.Temp = temp
	}

	return info, nil
}

func cpuUsage() (float64, error) {
	data, err := readProc("stat")
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

func cpuTemp() (float64, error) {
	matches, err := filepath.Glob("/sys/class/thermal/thermal_zone*/temp")
	if err != nil || len(matches) == 0 {
		return 0, fmt.Errorf("no thermal zones")
	}

	for _, m := range matches {
		base := strings.TrimSuffix(m, "/temp")
		typeData, _ := os.ReadFile(filepath.Join(base, "type"))
		typ := strings.TrimSpace(string(typeData))

		if strings.Contains(strings.ToLower(typ), "cpu") || strings.Contains(strings.ToLower(typ), "core") || strings.Contains(strings.ToLower(typ), "x86") {
			val, err := readUint64(m)
			if err != nil {
				continue
			}
			return float64(val) / 1000.0, nil
		}
	}

	val, err := readUint64(matches[0])
	if err != nil {
		return 0, err
	}
	return float64(val) / 1000.0, nil
}

func readMemory() (*MemoryInfo, error) {
	data, err := readProc("meminfo")
	if err != nil {
		return nil, err
	}

	info := &MemoryInfo{}
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

type mountEntry struct {
	device     string
	mountPoint string
	fstype     string
}

var realFS = map[string]bool{
	"ext4": true, "ext3": true, "ext2": true,
	"xfs": true, "btrfs": true, "zfs": true,
	"ntfs": true, "vfat": true, "exfat": true,
	"f2fs": true, "reiserfs": true, "jfs": true,
	"nfs": true, "nfs4": true, "cifs": true,
	"overlay": true,
}

func readMounts() ([]mountEntry, error) {
	data, err := readProc("mounts")
	if err != nil {
		return nil, err
	}

	var entries []mountEntry
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		fstype := fields[2]
		if !realFS[fstype] && !strings.HasPrefix(fstype, "fuse") {
			continue
		}

		mount := fields[1]
		device := fields[0]

		entries = append(entries, mountEntry{
			device:     device,
			mountPoint: mount,
			fstype:     fstype,
		})
	}

	return entries, nil
}

func statUsage(path string) (total, used, avail uint64, err error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, 0, 0, err
	}

	total = stat.Blocks * uint64(stat.Bsize)
	avail = stat.Bavail * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	used = total - free

	return total, used, avail, nil
}

func readDisk() (*DiskInfo, error) {
	mounts, err := readMounts()
	if err != nil {
		return nil, err
	}

	info := &DiskInfo{}
	for _, m := range mounts {
		total, used, avail, err := statUsage(m.mountPoint)
		if err != nil || total == 0 {
			continue
		}

		usage := float64(used) / float64(total) * 100

		info.Mounts = append(info.Mounts, MountInfo{
			Filesystem: m.fstype,
			MountPoint: m.mountPoint,
			Size:       total,
			Used:       used,
			Available:  avail,
			Usage:      usage,
		})
	}

	return info, nil
}

func readGPU() (*GPUInfo, error) {
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

		vendor := mapVendor(vendorID)
		if vendor == "" {
			continue
		}

		devDir := filepath.Dir(vPath)
		vram, _ := readUint64(filepath.Join(devDir, "mem_info_vram_total"))
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

	if usage, err := readInt(filepath.Join(best.dir, "gpu_busy_percent")); err == nil {
		info.Usage = float64(usage)
	}

	if total, err := readUint64(filepath.Join(best.dir, "mem_info_vram_total")); err == nil {
		if used, err := readUint64(filepath.Join(best.dir, "mem_info_vram_used")); err == nil {
			info.Memory = used
			_ = total
		}
	}

	info.Temp = gpuTemp()
	return info, nil
}

func gpuTemp() float64 {
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
		val, err := readInt(filepath.Join(hwmonDir, "temp1_input"))
		if err != nil {
			return 0
		}
		return float64(val) / 1000.0
	}

	return 0
}

func mapVendor(id string) string {
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

func readTemps() (*TempInfo, error) {
	info := &TempInfo{}

	sensors, err := readHwmonTemps()
	if err == nil {
		info.Sensors = append(info.Sensors, sensors...)
	}

	thermal, err := readThermalTemps()
	if err == nil {
		info.Sensors = append(info.Sensors, thermal...)
	}

	if len(info.Sensors) == 0 {
		return nil, fmt.Errorf("no temperature sensors found")
	}

	return info, nil
}

func readHwmonTemps() ([]TempSensor, error) {
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
			val, err := readInt(inputPath)
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

func readThermalTemps() ([]TempSensor, error) {
	matches, err := filepath.Glob("/sys/class/thermal/thermal_zone*")
	if err != nil {
		return nil, err
	}

	var zones []TempSensor
	for _, dir := range matches {
		val, err := readInt(filepath.Join(dir, "temp"))
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
			rxBytes, _ = readUint64(filepath.Join(ifPath, "statistics", "rx_bytes"))
			txBytes, _ = readUint64(filepath.Join(ifPath, "statistics", "tx_bytes"))
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
		if capVal, err := readInt(filepath.Join(dir, "capacity")); err == nil {
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

func readSystem() (*SystemInfo, error) {
	info := &SystemInfo{}

	hostname, err := os.Hostname()
	if err == nil {
		info.Hostname = hostname
	}

	if data, err := readProc("version"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) >= 3 {
			info.Kernel = fields[2]
		}
	}

	if data, err := readProc("uptime"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) > 0 {
			if sec, err := strconv.ParseFloat(fields[0], 64); err == nil {
				info.UptimeSec = sec
				info.Uptime = formatUptime(sec)
			}
		}
	}

	if data, err := readProc("loadavg"); err == nil {
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
