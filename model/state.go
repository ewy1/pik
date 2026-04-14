package model

type State struct {
	Sources []*Source
}

type HydratedState struct {
	*State
	HydratedSources []*HydratedSource
}
