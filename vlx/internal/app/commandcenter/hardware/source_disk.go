package hardware

import (
	"context"
	"fmt"

	"github.com/rlvelte/velinux/vlx/internal/core/printer"
)

func init() {
	register(&diskSource{})
}

type diskSource struct{}

func (s *diskSource) Name() string      { return "disk" }
func (s *diskSource) Aliases() []string { return []string{"disks", "storage", "df"} }

func (s *diskSource) Run(ctx context.Context, p *printer.Printer, json bool) error {
	info, err := readDisk()
	if err != nil {
		return fmt.Errorf("reading disk info: %w", err)
	}

	if json {
		PrintJSON(info)
		return nil
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
			formatBytes(m.Size),
			formatBytes(m.Used),
			formatBytes(m.Available),
			fmt.Sprintf("%.1f%%", m.Usage),
		})
	}

	p.Table([]string{"Type", "Mount", "Size", "Used", "Avail", "Use%"}, rows)
	return nil
}
