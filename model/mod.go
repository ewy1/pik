package model

type Modder interface {
	Mod(source *Source, result *HydratedSource) error
}
