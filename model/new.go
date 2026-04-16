package model

import (
	"errors"
	"io/fs"
	"path/filepath"
	"pik/identity"
	"strings"
	"sync"
)

func NewState(f fs.FS, locations []string, indexers []Indexer, runners []Runner) (*State, []error) {
	var errs []error
	st := &State{}
	wg := sync.WaitGroup{}
	for _, loc := range locations {
		wg.Go(func() {
			_, dirName := filepath.Split(loc)
			src := &Source{
				Path:     loc,
				Identity: identity.New(dirName),
			}
			loc = strings.TrimSuffix(loc, "/")
			loc = strings.TrimPrefix(loc, "/")

			if loc == "" {
				return
			}

			myWg := sync.WaitGroup{}
			for _, indexer := range indexers {
				myWg.Go(func() {
					s, err := fs.Sub(f, loc)
					if err != nil && !errors.Is(err, fs.ErrNotExist) {
						errs = append(errs, err)
						return
					}
					targets, err := indexer.Index("/"+loc, s, runners)
					if err != nil && !errors.Is(err, fs.ErrNotExist) {
						errs = append(errs, err)
						return
					}
					src.Targets = append(src.Targets, targets...)
				})
			}
			myWg.Wait()

			if src.Targets != nil {
				st.Sources = append(st.Sources, src)
			}
		})

	}
	wg.Wait()

	return st, errs
}
