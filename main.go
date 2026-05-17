package main

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/ewy1/pik/cache"
	"github.com/ewy1/pik/crawl"
	"github.com/ewy1/pik/flags"
	"github.com/ewy1/pik/git"
	"github.com/ewy1/pik/help"
	"github.com/ewy1/pik/indexers/pikdex"
	"github.com/ewy1/pik/menu"
	"github.com/ewy1/pik/model"
	"github.com/ewy1/pik/paths"
	"github.com/ewy1/pik/run"
	"github.com/ewy1/pik/runner/exc"
	"github.com/ewy1/pik/runner/gnumake"
	"github.com/ewy1/pik/runner/js"
	"github.com/ewy1/pik/runner/just"
	"github.com/ewy1/pik/runner/python"
	"github.com/ewy1/pik/runner/shell"
	"github.com/ewy1/pik/search"
	"github.com/ewy1/pik/spool"
	"github.com/spf13/pflag"
	"os"
	"runtime/pprof"
)

// syncInitializers are ran before the initializers.
// useful for initializing stuff like paths, preparing directories, and reading the environment
var syncInitializers = ComponentList[model.Initializer]{
	paths.Component,
}

// initializers are ran before indexing with the indexers,
// data from the syncInitializers can be accessed at this time.
var initializers = ComponentList[model.Initializer]{
	pikdex.Indexer,
	python.Python,
	git.Git,
	js.Js,
}

// indexers are methods which scan a directory and return a number of targets.
var indexers = ComponentList[model.Indexer]{
	pikdex.Indexer,
	just.Indexer,
	gnumake.Indexer,
	js.Js,
}

// runners are modules which know how to turn a file into an exec.Cmd
// all indexers have access to these but only pikdex uses it
var runners = ComponentList[model.Runner]{
	shell.Runner,
	python.Python,
	exc.Exc,
	js.Js,
}

// hydrators are ran when the menu is required
// for example adding git info, descriptions, icons...
var hydrators = ComponentList[model.Modder]{
	pikdex.Indexer,
	git.Git,
}

// ForceConfirm means we will have to ask for confirmation before running no matter what
var ForceConfirm = false

// SourcesWithoutResults is a failed cache from the previous iteration
// used for stripping out results to prevent double-index
var SourcesWithoutResults *cache.Cache

//go:embed version.txt
var version string

func main() {
	result := pik()
	if profileFd != nil {
		pprof.StopCPUProfile()
		err := profileFd.Close()
		if err != nil {
			_, _ = spool.Warn("%v\n", err)
		}
	}
	if result != 0 {
		os.Exit(result)
	}
}

func mode[T any](list ModeMap[T], fire func(mode T) error) *int {
	err := list.Traverse(func(in T) error {
		return fire(in)
	})
	if errors.Is(err, Success) {
		zero := 0
		return &zero
	} else if err != nil {
		_, _ = spool.Warn("%v\n", err)
		one := 1
		return &one
	}
	return nil
}

func pik() int {
	pflag.Usage = help.Echo
	pflag.Parse()

	code := mode(uninitializedModes, func(mode func() error) error {
		return mode()
	})
	if code != nil {
		return *code
	}

	code = mode(statelessModes, func(mode func() error) error {
		return mode()
	})
	if code != nil {
		return *code
	}

	syncInitializers.RunSync(func(initializer model.Initializer) error {
		return initializer.Init()
	})
	initializers.RunAsync(func(initializer model.Initializer) error {
		return initializer.Init()
	})

	here, err := os.Getwd()
	if err != nil {
		_, _ = spool.Warn("%v\n", err)
		return 1
	}
	locs := append(crawl.RichLocations(here), paths.System.String())
	root, err := os.OpenRoot("/")
	defer root.Close()
	if root == nil {
		_, _ = spool.Warn("%v\n", err)
		return 1
	}
	fs := root.FS()
	if err != nil {
		_, _ = spool.Warn("%v\n", err)
		return 1
	}
	var st *model.State
	var stateErrors []error

	var c *cache.Cache
	if !*flags.All {
		st, stateErrors = model.NewState(fs, locs, indexers, runners)
		go func() {
			if len(stateErrors) > 0 {
				return
			}
			err := cache.MergeAndSave(st)
			if err != nil {
				_, _ = spool.Warn("%v\n", err)
			}

		}()
	} else {
		c, err = cache.LoadFile(fs, paths.ContextsFile.String()[1:])
		if err != nil {
			_, _ = spool.Warn("%v\n", err)
			return 1
		}
		st, stateErrors = cache.LoadState(fs, c, indexers, runners)
	}
	if stateErrors != nil {
		_, _ = spool.Warn("%v\n", stateErrors)
	}

	code = mode(statefulModes, func(mode func(st *model.State) error) error {
		return mode(st)
	})
	if code != nil {
		return *code
	}

	args := pflag.Args()

	var result *search.Result
	cancelled := false

	if len(args) == 0 {
		md, err := menu.Show(st, hydrators)
		if err != nil {
			_, _ = spool.Warn("%v\n", err)
			return 1
		}
		cancelled = md.Cancel
		source, target := md.Result()
		if target != nil {
			t := target.Target()
			result = &search.Result{
				Target:            t,
				Source:            source.Source,
				NeedsConfirmation: false,
				Overridden:        t.Tags().Has(model.Override),
				Sub:               t.Sub(),
			}
		}
	}

	if result == nil {
		result = search.Search(st, args...)
	}

	// TODO: Move auto-all logic into Search?
	if !*flags.All && result.Target == nil && len(result.Args) > 0 && SourcesWithoutResults == nil && !ForceConfirm {
		ForceConfirm = true
		if err != nil {
			_, _ = spool.Warn("%v\n", err)
			return 1
		}
		SourcesWithoutResults = c
		return pik()
	}

	if cancelled {
		_, _ = spool.Warn("no target selected\n")
		return 0
	}

	if result.Target == nil {
		_, _ = spool.Warn("target not found\n")
		return 1
	}

	if result.NeedsConfirmation || ForceConfirm {
		_, _ = fmt.Fprintf(os.Stderr, "this target is out of tree.\n")
		if !menu.Confirm(os.Stdin, result.Source, result.Target, args...) {
			return 0
		}
	}
	if result.Overridden {
		_, _ = fmt.Fprintln(os.Stderr, menu.OverrideWarning(result.Target))
	}

	code = mode(selectionModes, func(mode func(st *model.State, src *model.Source, t model.Target) error) error {
		return mode(st, result.Source, result.Target)
	})
	if code != nil {
		return *code
	}

	err = run.Run(result.Source, result.Target, result.Args...)

	if err != nil {
		_, _ = spool.Warn("%v\n", err)
		return 1
	}

	return 0
}
