package model

type Hydrator interface {
	Hydrate(source *Source, result *HydratedSource) error
}
