package sources

type Source interface {
	Name() string
	CheckUpdates() (count int, packages []string, err error)
}
