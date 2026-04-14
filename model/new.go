package model

import (
	"io/fs"
	"path/filepath"
	"pik/identity"
	"strings"
)

func NewState(f fs.FS, locations []string, indexers []Indexer, runners []Runner) (*State, error) {
	st := &State{}
	for _, loc := range locations {
		_, dirName := filepath.Split(loc)
		src := &Source{
			Path:     loc,
			Identity: identity.New(dirName),
		}
		loc = strings.TrimSuffix(loc, "/")
		loc = strings.TrimPrefix(loc, "/")

		if loc == "" {
			continue
		}

		for _, indexer := range indexers {

			s, err := fs.Sub(f, loc)
			if err != nil {
				return nil, err
			}
			targets, err := indexer.Index("/"+loc, s, runners)
			if err != nil {
				return nil, err
			}
			src.Targets = append(src.Targets, targets...)
		}

		if src.Targets != nil {
			st.Sources = append(st.Sources, src)
		}

	}

	return st, nil
}
