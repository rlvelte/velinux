package sources

import (
	"fmt"
	"strings"

	"github.com/rlvelte/velinux/vlx/internal/visuals/printer"
)

// sources compiles a list of all available data readings.
var sources = []Source{
	&cpu{},
	&disk{},
	&gpu{},
	&network{},
	&power{},
	&ram{},
	&temperature{},
	&system{},
}

// Source is a data provider for hardware info.
type Source interface {
	Name() string
	Aliases() []string
	Run(p *printer.Printer) error
}

// Find looks up a source by name or alias.
func Find(name string) (Source, error) {
	for _, s := range sources {
		if s.Name() == name {
			return s, nil
		}

		for _, a := range s.Aliases() {
			if a == name {
				return s, nil
			}
		}
	}

	return nil, fmt.Errorf("unknown hardware source: %s", name)
}

// List prints all registered sources.
func List(p *printer.Printer) {
	rows := make([][]string, 0, len(sources))
	for _, s := range sources {
		aliases := ""
		if len(s.Aliases()) > 0 {
			aliases = strings.Join(s.Aliases(), ", ")
		}

		rows = append(rows, []string{s.Name(), aliases})
	}

	p.Table([]string{"Source", "Aliases"}, rows)
}

func FormatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}

	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
