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
	"runtime"
	"runtime/pprof"
)

// ModeMap maps flags to specific operation modes
type ModeMap[T any] map[*bool]T

// Continue can be returned as an error to continue program flow
var Continue = errors.New("not an error; continue flow")
var Success = errors.New("not an error; finished operations")

// Traverse checks the entries of the map. If any flags are set on,
// pik that mode. If Continue is returned, it's non-exclusive. Otherwise,
// we quit after one mode.
//
// `then` should simply be the method call (necessary due to generics)
func (m ModeMap[T]) Traverse(then func(in T) error) error {
	for enabled, mode := range m {
		if !*enabled {
			continue
		}
		err := then(mode)
		if errors.Is(err, Continue) {
			continue
		} else if err != nil {
			return err
		}
		return Success
	}
	return nil
}

var profileFd *os.File

// statelessModes are program modes which do not require state to operate.
// like --version and --completion
var statelessModes = ModeMap[func() error]{
	flags.Version: func() error {
		_, err := spool.Print("%s\n", version)
		return err
	},
	flags.Completion: func() error {
		return completion.Echo()
	},
	flags.Profile: func() error {
		fd, err := os.Create("pik-profile.out")
		if err != nil {
			return err
		}
		runtime.SetCPUProfileRate(1000)
		err = pprof.StartCPUProfile(profileFd)
		if err != nil {
			return err
		}
		profileFd = fd
		return Continue
	},
}

// statefulModes are program modes which require a built state to be executed
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

// selectionModes are program modes which require a selected target, through menu or args
var selectionModes = ModeMap[func(st *model.State, src *model.Source, t model.Target) error]{
	flags.Edit: func(st *model.State, src *model.Source, t model.Target) error {
		return run.Edit(t, src)
	},
}
