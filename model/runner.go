package model

import (
	"io/fs"
)

type Runner interface {
	Hydrate(target Target) (HydratedTarget, error)
	Wants(fs fs.FS, file string, entry fs.DirEntry) (bool, error)
	CreateTarget(fs fs.FS, source string, file string, entry fs.DirEntry) (Target, error)
}
