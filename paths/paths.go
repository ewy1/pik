package paths

import (
	"github.com/adrg/xdg"
	"os"
	"path/filepath"
	"strings"
)

var Home, This, Cache, Config, Ifs string

type paths struct {
	Initialized bool
}

var Paths = &paths{
	Initialized: false,
}

func (p *paths) Init() error {

	// if we're asked to initialize for a second time,
	// probably some environment has changed
	if p.Initialized {
		xdg.Reload()
	}

	Home = xdg.Home
	This = "pik"
	Cache = filepath.Join(xdg.CacheHome, This)
	Config = filepath.Join(xdg.ConfigHome, This)
	Ifs = os.Getenv("IFS")

	err := os.MkdirAll(Cache, 0700)
	if err != nil {
		return err
	}
	err = os.MkdirAll(Config, 0700)
	if err != nil {
		return err
	}
	if Ifs == "" {
		Ifs = "\n"
	}
	p.Initialized = true
	return nil
}

func ReplaceHome(input string) string {
	return strings.Replace(input, Home, "~", 1)
}
