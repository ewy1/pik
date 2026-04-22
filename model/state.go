package model

type State struct {
	Sources []*Source
	All     bool
}

type HydratedState struct {
	*State
	HydratedSources []*HydratedSource
}
