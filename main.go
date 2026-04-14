package main

import (
	"github.com/spf13/pflag"
	"os"
	"pik/cache"
	"pik/crawl"
	"pik/flags"
	"pik/indexers/pikdex"
	"pik/menu"
	"pik/model"
	"pik/run"
	"pik/runner/just"
	"pik/runner/python"
	"pik/runner/shell"
	"pik/search"
	"pik/spool"
)

var initializers = []model.HasInit{
	python.Python,
}

var indexers = []model.Indexer{
	pikdex.Indexer,
	just.Indexer,
}

var runners = []model.Runner{
	shell.Runner,
	python.Python,
}

var hydrators = []model.Hydrator{
	pikdex.Indexer,
}

var ForceConfirm = false

func main() {
	pflag.Parse()
	for _, i := range initializers {
		err := i.Init()
		if err != nil {
			_, _ = spool.Warn("%v\n", err)
		}
	}
	here, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	locs := crawl.RichLocations(here)
	last := locs[len(locs)-1]
	root, err := os.OpenRoot(last)
	if root == nil {
		panic(err)
	}
	fs := root.FS()
	if err != nil {
		panic(err)
	}
	var st *model.State
	var stateErrors []error
	if !*flags.All {
		st, stateErrors = model.NewState(fs, locs, indexers, runners)
	} else {
		c, err := cache.Load()
		if err != nil {
			panic(err)
		}
		st, stateErrors = cache.LoadState(fs, c, indexers, runners)
	}
	if stateErrors != nil {
		spool.Warn("%v\n", stateErrors)
	} else {
		err = cache.Save(st)
	}

	if err != nil {
		_, _ = spool.Warn("%v", err)
	}

	args := pflag.Args()

	if len(args) == 0 {
		source, target, err := menu.Show(st, hydrators)
		if err != nil {
			spool.Warn("%v\n", err)
			os.Exit(1)
		}
		if target == nil {
			spool.Warn("no target selected.\n")
			os.Exit(0)
		}
		err = run.Run(source.Source, target, args...)
		if err != nil {
			spool.Warn("%v\n", err)
			os.Exit(1)
		}

		return
	}

	target, src, confirm, _, args := search.Search(st, args...)
	if !*flags.All && target == nil && len(args) > 0 {
		err := pflag.Set("all", "true")
		ForceConfirm = true
		if err != nil {
			spool.Warn("%v\n", err)
			os.Exit(1)
		}
		main()
		return
	}

	if target == nil {
		_, _ = spool.Print("no target found.")
		os.Exit(1)
		return
	}

	if confirm || ForceConfirm {
		if !menu.Confirm(os.Stdin, src, target, args...) {
			os.Exit(0)
		}
	}

	err = run.Run(src, target, args...)
	if err != nil {
		spool.Warn("%v\n", err)
		os.Exit(1)
	}
}
