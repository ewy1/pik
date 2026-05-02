package main

import (
	"errors"
	"github.com/ewy1/pik/completion"
	"github.com/ewy1/pik/flags"
	"github.com/ewy1/pik/model"
	"github.com/ewy1/pik/paths"
	"github.com/ewy1/pik/run"
	"github.com/ewy1/pik/spool"
	"os"
)

type ModeMap[T any] map[*bool]T

func (m ModeMap[T]) Traverse(then func(in T) error) {
	for enabled, mode := range m {
		if !*enabled {
			continue
		}
		err := then(mode)
		if errors.Is(err, Continue) {
			continue
		} else if err != nil {
			_, _ = spool.Warn("%v\n", err)
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}
}

var statelessModes = ModeMap[func() error]{
	flags.Version: func() error {
		_, err := spool.Print("%s\n", version)
		return err
	},
	flags.Completion: func() error {
		return completion.Echo()

	},
}

var Continue = errors.New("not an error; continue flow")

var statefulModes = ModeMap[func(st *model.State) error]{
	flags.List: func(st *model.State) error {
		for _, s := range st.Sources {
			_, _ = spool.Print("%v", s.Label()+paths.Ifs)
			for _, t := range s.Targets {
				_, _ = spool.Print("%v", t.ShortestId()+paths.Ifs)
			}
		}
		return nil
	},
}

var selectionModes = ModeMap[func(st *model.State, src *model.Source, t model.Target) error]{
	flags.Edit: func(st *model.State, src *model.Source, t model.Target) error {
		return run.Edit(t, src)
	},
}
