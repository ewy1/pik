package main

import (
	_ "embed"
	"fmt"
	"github.com/ewy1/pik/cache"
	"github.com/ewy1/pik/crawl"
	"github.com/ewy1/pik/flags"
	"github.com/ewy1/pik/git"
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
)

// syncInitializers are ran before the initializers.
// useful for initializing stuff like paths, preparing directories, and reading the environment
var syncInitializers = ComponentList[model.Initializer]{
	paths.Paths,
	cache.Init,
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
	pflag.Parse()

	statelessModes.Traverse(func(in func() error) error {
		return in()
	})

	syncInitializers.RunSync(func(initializer model.Initializer) error {
		return initializer.Init()
	})

	initializers.RunAsync(func(initializer model.Initializer) error {
		return initializer.Init()
	})

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

	var c *cache.Cache
	if !*flags.All {
		st, stateErrors = model.NewState(fs, locs, indexers, runners)
		err = cache.Insert(st)
		if err != nil {
			spool.Warn("%v\n", err)
		}
	} else {
		c, err = cache.LoadFile(fs, cache.Path[1:])
		if err != nil {
			_, _ = spool.Warn("%v\n", err)
			os.Exit(1)
		}
		st, stateErrors = cache.LoadState(fs, c, indexers, runners)
	}
	if stateErrors != nil {
		_, _ = spool.Warn("%v\n", stateErrors)
	}

	statefulModes.Traverse(func(in func(st *model.State) error) error {
		return in(st)
	})

	args := pflag.Args()

	var result *search.Result
	cancelled := false

	if len(args) == 0 {
		md, err := menu.Show(st, hydrators)
		if err != nil {
			_, _ = spool.Warn("%v\n", err)
			os.Exit(1)
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
			os.Exit(1)
		}
		SourcesWithoutResults = c
		main()
		return
	}

	if cancelled {
		_, _ = spool.Warn("no target selected\n")
		os.Exit(0)
		return
	}

	if result.Target == nil {
		_, _ = spool.Warn("target not found\n")
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

	selectionModes.Traverse(func(in func(st *model.State, src *model.Source, t model.Target) error) error {
		return in(st, result.Source, result.Target)
	})

	err = run.Run(result.Source, result.Target, result.Args...)

	if err != nil {
		_, _ = spool.Warn("%v\n", err)
		os.Exit(1)
	}
}
