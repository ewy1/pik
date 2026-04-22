package model

import (
	"os/exec"
)

// Target is something we want to run
type Target interface {
	Matches
	// Create creates the exec.Cmd from itself. The model.Source this target is in is also available.
	Create(s *Source) *exec.Cmd
	Sub() []string
	// Label will be used for its label in the menu
	Label() string
	Hydrate(src *Source) (HydratedTarget, error)
	Tags() Tags
	// ShortestId returns a short version of its identity
	ShortestId() string
	// Visible should return whether this target can be seen in the menu
	Visible() bool
	// Invocation should return the "canonical invocation": simple to remember
	Invocation(src *Source) []string
}

// HydratedTarget is something we want to show in the menu
type HydratedTarget interface {
	Target
	// Icon is some text which will be used as an icon
	Icon() string
	// Description is a one-line description of what this does
	Description() string
	// Target returns our inner target
	Target() Target
}
