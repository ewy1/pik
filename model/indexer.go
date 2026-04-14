package model

import (
	"io/fs"
)

type Indexer interface {
	Index(path string, fs fs.FS, runners []Runner) ([]Target, error)
}
