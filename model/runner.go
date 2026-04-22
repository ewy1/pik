package model

import (
	"io/fs"
)

// Runner knows how to go from a file to an exec.Cmd
// these are mostly used for pikdex but other runners can use them as well
// these have to be registered in main.go
type Runner interface {
	Hydrate(target Target) (HydratedTarget, error)
	Wants(fs fs.FS, file string, entry fs.DirEntry) (bool, error)
	CreateTarget(fs fs.FS, source string, file string, entry fs.DirEntry) (Target, error)
}
