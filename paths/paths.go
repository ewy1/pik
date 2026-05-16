package paths

import (
	"github.com/adrg/xdg"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	Ifs          string
	ConfigDir    = Empty()
	CacheDir     = Empty()
	ContextsFile = Empty()
	HomeDir      = Empty()
	System       = Empty()
	This         = "pik"
)

var Paths = []*Path{
	ConfigDir,
	CacheDir,
	ContextsFile,
	HomeDir,
	System,
}

var DirectoriesToMake = []*Path{
	ConfigDir,
	CacheDir,
}

type paths struct {
	Initialized bool
}

var Component = &paths{
	Initialized: false,
}

func (p *paths) Init() error {

	// if we're asked to initialize for a second time,
	// probably some environment has changed
	if p.Initialized {
		xdg.Reload()
	}

	HomeDir.Set(xdg.Home)
	CacheDir.Set(filepath.Join(xdg.CacheHome, This))
	ConfigDir.Set(filepath.Join(xdg.ConfigHome, This))
	ContextsFile.Set(path.Join(CacheDir.String(), "contexts"))
	System.Set(path.Join(ConfigDir.String()))

	Ifs = os.Getenv("IFS")

	for _, d := range DirectoriesToMake {
		err := os.MkdirAll(d.String(), 0700)
		if err != nil {
			return err
		}
	}

	if Ifs == "" {
		Ifs = "\n"
	}
	p.Initialized = true
	return nil
}

func ReplaceHome(input string) string {
	return strings.Replace(input, HomeDir.String(), "~", 1)
}
