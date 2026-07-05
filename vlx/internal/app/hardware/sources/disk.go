package sources

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"syscall"

	"github.com/rlvelte/velinux/vlx/internal/core/format"
	"github.com/rlvelte/velinux/vlx/internal/core/fsys"
	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
)

// DiskInfo contains information about disk.
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

// disk information container.
type disk struct{}

// Name of this source.
func (s *disk) Name() string {
	return "disk"
}

// Aliases that this source has.
func (s *disk) Aliases() []string {
	return []string{"disks", "storage", "df"}
}

// Run extracts all data for this source.
func (s *disk) Run(p *printer.Printer) error {
	info, err := readDisk()
	if err != nil {
		return fmt.Errorf("reading disk info: %w", err)
	}

	if len(info.Mounts) == 0 {
		p.Info("No mounted filesystems found")
		return nil
	}

	rows := make([][]string, 0, len(info.Mounts))
	for _, m := range info.Mounts {
		rows = append(rows, []string{
			m.Filesystem,
			m.MountPoint,
			format.Bytes(m.Size),
			format.Bytes(m.Used),
			format.Bytes(m.Available),
			fmt.Sprintf("%.1f%%", m.Usage),
		})
	}

	p.Table([]string{"Type", "Mount", "Size", "Used", "Avail", "Use%"}, rows)
	return nil
}

// readDisk aggregates all available data for disk.
func readDisk() (*DiskInfo, error) {
	mounts, err := readDiskMounts()
	if err != nil {
		return nil, err
	}

	info := &DiskInfo{}
	for _, m := range mounts {
		total, used, avail, err := readDiskUsage(m.mountPoint)
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

// readDiskMounts aggregates mount entries for disk.
func readDiskMounts() ([]mountEntry, error) {
	data, err := fsys.Proc("mounts")
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

// readDiskUsage aggregates usage data for disk.
func readDiskUsage(path string) (total, used, avail uint64, err error) {
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
