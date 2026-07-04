package sources

import (
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
func Find(name string) Source {
	for _, s := range sources {
		if s.Name() == name {
			return s
		}

		for _, a := range s.Aliases() {
			if a == name {
				return s
			}
		}
	}

	return nil
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
