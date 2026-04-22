package model

import (
	"errors"
	"io/fs"
	"path/filepath"
	"pik/flags"
	"pik/identity"
	"strings"
	"sync"
)

func NewState(f fs.FS, locations []string, indexers []Indexer, runners []Runner) (*State, []error) {
	var errs []error
	st := &State{
		All: *flags.All,
	}
	wg := sync.WaitGroup{}
	var sources = make([]*Source, len(locations), len(locations))
	for i, loc := range locations {
		wg.Go(func() {
			_, dirName := filepath.Split(strings.TrimSuffix(loc, "/"))
			src := &Source{
				Path:     loc,
				Identity: identity.New(dirName),
			}
			sources[i] = src
			loc = strings.TrimSuffix(loc, "/")
			loc = strings.TrimPrefix(loc, "/")

			if loc == "" {
				return
			}

			myWg := sync.WaitGroup{}
			var targets = make([][]Target, len(indexers), len(indexers))
			for ti, indexer := range indexers {
				myWg.Go(func() {
					s, err := fs.Sub(f, loc)
					if err != nil && !errors.Is(err, fs.ErrNotExist) {
						errs = append(errs, err)
						return
					}
					result, err := indexer.Index("/"+loc, s, runners)
					if err != nil && !errors.Is(err, fs.ErrNotExist) {
						errs = append(errs, err)
						return
					}
					targets[ti] = result
				})
			}
			myWg.Wait()

			for _, t := range targets {
				if t == nil {
					continue
				}
				sources[i].Targets = append(sources[i].Targets, t...)
			}

		})

	}
	wg.Wait()

	for _, s := range sources {
		if s == nil || s.Targets == nil {
			continue
		}
		st.Sources = append(st.Sources, s)
	}

	return st, errs
}
