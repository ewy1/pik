package model

import (
	"os/exec"
)

type Target interface {
	Matches
	Create(s *Source) *exec.Cmd
	Sub() []string
	Label() string
	Hydrate(src *Source) (HydratedTarget, error)
	Tags() Tags
	ShortestId() string
	Visible() bool
	Invocation(src *Source) []string
}

type HydratedTarget interface {
	Target
	Icon() string
	Description() string
	Target() Target
}
