package model

import (
	"io/fs"
)

// Indexer is used to create targets from pik-enabled sources
// these have to be listed in main.go
type Indexer interface {
	Index(path string, fs fs.FS, runners []Runner) ([]Target, error)
}
