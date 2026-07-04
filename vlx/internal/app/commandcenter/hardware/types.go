package hardware

type CPUInfo struct {
	Model    string  `json:"model"`
	Cores    int     `json:"cores"`
	Threads  int     `json:"threads"`
	Freq     float64 `json:"freq_mhz"`
	MinFreq  float64 `json:"min_freq_mhz,omitempty"`
	MaxFreq  float64 `json:"max_freq_mhz,omitempty"`
	Usage    float64 `json:"usage_percent"`
	Temp     float64 `json:"temp_celsius,omitempty"`
}

type MemoryInfo struct {
	Total     uint64 `json:"total_bytes"`
	Used      uint64 `json:"used_bytes"`
	Free      uint64 `json:"free_bytes"`
	Available uint64 `json:"available_bytes"`
	SwapTotal uint64 `json:"swap_total_bytes"`
	SwapUsed  uint64 `json:"swap_used_bytes"`
	SwapFree  uint64 `json:"swap_free_bytes"`
}

type DiskInfo struct {
	Mounts []MountInfo `json:"mounts"`
}

type MountInfo struct {
	Filesystem string  `json:"filesystem"`
	MountPoint string  `json:"mount_point"`
	Size       uint64  `json:"size_bytes"`
	Used       uint64  `json:"used_bytes"`
	Available  uint64  `json:"available_bytes"`
	Usage      float64 `json:"usage_percent"`
}

type GPUInfo struct {
	Vendor  string  `json:"vendor"`
	Model   string  `json:"model"`
	Driver  string  `json:"driver,omitempty"`
	Memory  uint64  `json:"memory_bytes,omitempty"`
	Temp    float64 `json:"temp_celsius,omitempty"`
	Usage   float64 `json:"usage_percent,omitempty"`
}

type TempInfo struct {
	Sensors []TempSensor `json:"sensors"`
}

type TempSensor struct {
	Name  string  `json:"name"`
	Label string  `json:"label,omitempty"`
	Temp  float64 `json:"temp_celsius"`
}

type NetworkInfo struct {
	Interfaces []NetInterface `json:"interfaces"`
}

type NetInterface struct {
	Name    string `json:"name"`
	Up      bool   `json:"up"`
	RxBytes uint64 `json:"rx_bytes"`
	TxBytes uint64 `json:"tx_bytes"`
}

type PowerInfo struct {
	Supplies []PowerSupply `json:"supplies"`
}

type PowerSupply struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Status   string  `json:"status"`
	Capacity float64 `json:"capacity_percent,omitempty"`
}

type SystemInfo struct {
	Hostname  string      `json:"hostname"`
	OS        string      `json:"os"`
	Kernel    string      `json:"kernel"`
	Uptime    string      `json:"uptime"`
	UptimeSec float64     `json:"uptime_seconds"`
	LoadAvg   [3]float64  `json:"load_average"`
	Processes int         `json:"processes"`
}
