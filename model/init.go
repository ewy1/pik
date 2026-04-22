package model

// Initializer has setup to do before pik can start
// these have to be registered in main.go
type Initializer interface {
	Init() error
}
