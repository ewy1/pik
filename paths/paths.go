package paths

import (
	"github.com/adrg/xdg"
	"os"
	"path/filepath"
	"strings"
)

var (
	Home   = xdg.Home
	This   = "pik"
	Cache  = filepath.Join(xdg.CacheHome, This)
	Config = filepath.Join(xdg.ConfigHome, This)
	Ifs    = os.Getenv("IFS")
)

type paths struct {
}

var Paths = &paths{}

func (p paths) Init() error {
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
	return nil
}

func ReplaceHome(input string) string {
	return strings.Replace(input, Home, "~", 1)
}
