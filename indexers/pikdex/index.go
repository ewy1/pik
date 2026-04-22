package pikdex

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"pik/model"
	"pik/spool"
	"slices"
	"strings"
	"sync"
)

var Roots = []string{

	// current name
	".pik",
	"_pik",

	// program names from a previous life
	".godo",
	"_godo",
	".pik",
	"_uwu",

	//utility
	".bin",
	"_bin",
	"tasks",
	".tasks",
	"_tasks",
}

var SkippedFolders = []string{
	".git",
	".config",
	".idea",
}

func (u *pikdex) Init() error {
	// add own executable name to uwudexable dirs
	self, err := os.Executable()
	if strings.HasSuffix(self, ".test") {
		return nil
	}
	if err != nil {
		_, _ = spool.Warn("%v\n", err)
		return nil
	}
	self = strings.TrimSuffix(self, ".exe")
	Roots = append(Roots, "."+self, "_"+self)
	return nil
}

var Indexer = &pikdex{mods: make(map[string]*SourceData)}

type pikdex struct {
	sync.Mutex
	mods map[string]*SourceData
}

type SourceData struct {
	Aliases []string
	Icon    string
	Path    string
}

func (u *pikdex) Index(absPath string, f fs.FS, runners []model.Runner) ([]model.Target, error) {
	wants, root, err := u.WantsWalk(f)
	if !wants {
		return nil, err
	}
	var targets []model.Target
	mod := u.mods[absPath]
	if mod == nil {
		u.mods[absPath] = &SourceData{
			Path: absPath,
		}
		mod = u.mods[absPath]
	}
	err = fs.WalkDir(f, root, func(p string, d fs.DirEntry, err error) error {

		if !d.IsDir() {
			for trigger, applier := range MetaFiles {

				// during the crawl, we might find meta files
				expectedLocation := filepath.Join(absPath, root, trigger)
				actualLocation := filepath.Join(absPath, p)
				if expectedLocation != actualLocation {
					continue
				}

				content, err := os.ReadFile(expectedLocation)
				if err != nil {
					spool.Warn("%v\n", err)
					continue
				}
				applier(mod, string(content))

			}
		}

		if d.IsDir() {
			_, dirName := path.Split(p)
			if slices.Contains(SkippedFolders, dirName) {
				return fs.SkipDir
			}
		}

		for _, r := range runners {
			wants, err := r.Wants(f, p, d)
			if err != nil {
				spool.Warn("%v\n", err)
			}
			if wants {
				t, err := r.CreateTarget(f, absPath, p, d)
				if err != nil {
					spool.Warn("%v\n", err)
				}
				targets = append(targets, t)
				return nil
			}
			if err != nil {
				spool.Warn("%v\n", err)
			}
		}
		return nil
	})
	u.Lock()
	u.mods[absPath] = mod
	u.Unlock()

	return targets, err
}

func (u *pikdex) WantsWalk(f fs.FS) (bool, string, error) {
	entries, err := fs.ReadDir(f, ".")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, "", nil
		} else {
			return false, "", err
		}
	}

	for _, e := range entries {
		for _, r := range Roots {
			if e.Name() == r && e.IsDir() {
				return true, r, nil
			}
		}
	}

	return false, "", nil
}
