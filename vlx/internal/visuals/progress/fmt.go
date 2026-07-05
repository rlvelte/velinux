package progress

// FmtProgress is a fallback backend that uses plain fmt and os operations.
type FmtProgress struct{}

func (f *FmtProgress) Available() bool {
	return true
}

func (f *FmtProgress) Stream(message string) {
	panic("implement me")
}
