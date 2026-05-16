package model

import (
	"errors"
	"github.com/ewy1/pik/flags"
	"github.com/ewy1/pik/identity"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
)

func NewState(rootFs fs.FS, locations []string, indexers []Indexer, runners []Runner) (*State, []error) {
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

			locationWg := sync.WaitGroup{}
			var targets = make([][]Target, len(indexers), len(indexers))
			for ti, indexer := range indexers {
				locationWg.Go(func() {
					subFs, err := fs.Sub(rootFs, loc)
					if err != nil && !errors.Is(err, fs.ErrNotExist) {
						errs = append(errs, err)
						return
					}
					result, err := indexer.Index("/"+loc, subFs, runners)
					if err != nil && !errors.Is(err, fs.ErrNotExist) {
						errs = append(errs, err)
						return
					}
					targets[ti] = result
				})
			}
			locationWg.Wait()

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
