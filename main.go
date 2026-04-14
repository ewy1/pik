package main

import (
	_ "embed"
	"github.com/spf13/pflag"
	"os"
	"pik/cache"
	"pik/crawl"
	"pik/flags"
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
)

var initializers = []model.HasInit{
	python.Python,
}

var indexers = []model.Indexer{
	pikdex.Indexer,
	just.Indexer,
	gnumake.Indexer,
}

var runners = []model.Runner{
	shell.Runner,
	python.Python,
}

var hydrators = []model.Hydrator{
	pikdex.Indexer,
}

var ForceConfirm = false

//go:embed version.txt
var version string

func main() {
	pflag.Parse()

	switch {
	case *flags.Version:
		_, _ = spool.Print("%s\n", version)
		os.Exit(0)
	}

	for _, i := range initializers {
		err := i.Init()
		if err != nil {
			_, _ = spool.Warn("%v\n", err)
		}
	}
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

	if !*flags.All {
		st, stateErrors = model.NewState(fs, locs, indexers, runners)
	} else {
		c, err := cache.Load()
		if err != nil {
			_, _ = spool.Warn("%v\n", err)
			os.Exit(1)
		}
		st, stateErrors = cache.LoadState(fs, c, indexers, runners)
	}
	if stateErrors != nil {
		_, _ = spool.Warn("%v\n", stateErrors)
	} else {
		err = cache.Save(st)
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

	target, src, confirm, _, args := search.Search(st, args...)
	if !*flags.All && target == nil && len(args) > 0 {
		err := pflag.Set("all", "true")
		ForceConfirm = true
		if err != nil {
			_, _ = spool.Warn("%v\n", err)
			os.Exit(1)
		}
		main()
		return
	}

	if target == nil {
		_, _ = spool.Print("target not found.")
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
		_, _ = spool.Warn("%v\n", err)
		os.Exit(1)
	}
}
