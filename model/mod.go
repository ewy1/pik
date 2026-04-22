package model

// Modder can change sources after their creation
// this can add additional data, customization, etc.
// these have to be registered in main.go
type Modder interface {
	Mod(source *Source, result *HydratedSource) error
}
