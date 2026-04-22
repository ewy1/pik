package model

// State is the complete collection of sources and their targets
type State struct {
	Sources []*Source
	// All contains whether --all was enabled during creation
	All bool
}

// HydratedState is an extension of the State which contains
// additional information which we use to display the menu
type HydratedState struct {
	*State
	HydratedSources []*HydratedSource
}
