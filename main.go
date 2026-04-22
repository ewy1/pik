package main

import (
	_ "embed"
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"pik/cache"
	"pik/crawl"
	"pik/flags"
	"pik/git"
	"pik/indexers/pikdex"
	"pik/menu"
	"pik/model"
	"pik/paths"
	"pik/run"
	"pik/runner/gnumake"
	"pik/runner/just"
	"pik/runner/python"
	"pik/runner/shell"
	"pik/search"
	"pik/spool"
	"sync"
)

// preInitializers are ran before the initializers.
// useful for initializing stuff like paths, preparing directories, and reading the environment
var preInitializers = []model.Initializer{
	paths.Paths,
}

// initializers are ran before indexing with the indexers,
// data from the preInitializers can be accessed at this time.
var initializers = []model.Initializer{
	pikdex.Indexer,
	python.Python,
	git.Git,
}

// indexers are methods which scan a directory and return a number of targets.
var indexers = []model.Indexer{
	pikdex.Indexer,
	just.Indexer,
	gnumake.Indexer,
}

// runners are modules which know how to turn a file into an exec.Cmd
// all indexers have access to these but only pikdex uses it
var runners = []model.Runner{
	shell.Runner,
	python.Python,
}

// hydrators are ran when the menu is required
// for example adding git info, descriptions, icons...
var hydrators = []model.Modder{
	pikdex.Indexer,
	git.Git,
}

// ForceConfirm means we will have to ask for confirmation before running no matter what
var ForceConfirm = false

// SourcesWithoutResults is a failed cache from the previous iteration
// used for stripping out results to prevent double-index
var SourcesWithoutResults cache.Cache

//go:embed version.txt
var version string

func main() {
	pflag.Parse()

	switch {
	case *flags.Version:
		_, _ = spool.Print("%s\n", version)
		os.Exit(0)
	}

	wg := sync.WaitGroup{}
	for _, i := range preInitializers {
		wg.Go(func() {
			err := i.Init()
			if err != nil {
				_, _ = spool.Warn("%v\n", err)
			}
		})
	}
	wg.Wait()

	wg = sync.WaitGroup{}
	for _, i := range initializers {
		wg.Go(func() {
			err := i.Init()
			if err != nil {
				_, _ = spool.Warn("%v\n", err)
			}
		})
	}
	wg.Wait()

	here, err := os.Getwd()
	if err != nil {
		_, _ = spool.Warn("%v\n", err)
		os.Exit(1)
	}
	locs := crawl.RichLocations(here)
	last := locs[len(locs)-1]
	root, err := os.OpenRoot(last)
	if root == nil {
		_, _ = spool.Warn("%v\n", err)
		os.Exit(1)
	}
	fs := root.FS()
	if err != nil {
		_, _ = spool.Warn("%v\n", err)
		os.Exit(1)
	}
	var st *model.State
	var stateErrors []error

	var c cache.Cache
	if !*flags.All {
		st, stateErrors = model.NewState(fs, locs, indexers, runners)
	} else {
		c, err = cache.LoadFile(fs, cache.Path[1:])
		c.Strip(SourcesWithoutResults)
		if err != nil {
			_, _ = spool.Warn("%v\n", err)
			os.Exit(1)
		}
		st, stateErrors = cache.LoadState(fs, c, indexers, runners)
	}
	if stateErrors != nil {
		_, _ = spool.Warn("%v\n", stateErrors)
	} else {
		err = cache.SaveFile(cache.Path, st, c)
		if err != nil {
			_, _ = spool.Warn("%v", err)
		}
	}

	if *flags.List {
		for _, s := range st.Sources {
			for _, t := range s.Targets {
				_, _ = spool.Print(t.ShortestId() + paths.Ifs)
			}
		}
		os.Exit(0)
	}

	args := pflag.Args()

	if len(args) == 0 {
		source, target, err := menu.Show(st, hydrators)
		if err != nil {
			_, _ = spool.Warn("%v\n", err)
			os.Exit(1)
		}
		if target == nil {
			_, _ = spool.Warn("no target selected.\n")
			os.Exit(0)
		}
		err = run.Run(source.Source, target, args...)
		if err != nil {
			_, _ = spool.Warn("%v\n", err)
			os.Exit(1)
		}

		return
	}

	result := search.Search(st, args...)
	// TODO: Move auto-all logic into Search?
	if !*flags.All && result.Target == nil && len(args) > 0 {
		ForceConfirm = true
		if err != nil {
			_, _ = spool.Warn("%v\n", err)
			os.Exit(1)
		}
		SourcesWithoutResults = c
		main()
		return
	}

	if result.Target == nil {
		_, _ = spool.Print("target not found.")
		os.Exit(1)
		return
	}

	if result.NeedsConfirmation || ForceConfirm {
		_, _ = fmt.Fprintf(os.Stderr, "this target is out of tree.\n")
		if !menu.Confirm(os.Stdin, result.Source, result.Target, args...) {
			os.Exit(0)
		}
	}
	if result.Overridden {
		_, _ = fmt.Fprintln(os.Stderr, menu.OverrideWarning(result.Target))
	}
	err = run.Run(result.Source, result.Target, args...)
	if err != nil {
		_, _ = spool.Warn("%v\n", err)
		os.Exit(1)
	}
}
