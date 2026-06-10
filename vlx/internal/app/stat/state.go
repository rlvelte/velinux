package stat

import (
	"time"

	"github.com/rvelte/vlx/internal/system/state"
)

var statState = state.New("stat")

type SourceState struct {
	Count    int      `json:"count"`
	Packages []string `json:"packages"`
}

type UpdateState struct {
	LastCheck time.Time              `json:"last_check"`
	Sources   map[string]SourceState `json:"sources"`
}

func writeState(s *UpdateState) error {
	return statState.Write(s)
}

func readState() (*UpdateState, error) {
	s, err := state.Read[UpdateState](statState)
	if err != nil {
		return nil, err
	}

	if s.LastCheck.IsZero() && len(s.Sources) == 0 {
		return nil, nil
	}

	return &s, nil
}
