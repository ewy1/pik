package paths

import (
	"github.com/adrg/xdg"
	"os"
	"path/filepath"
	"pik/spool"
	"strings"
)

var (
	Home   = xdg.Home
	This   = "pik"
	Cache  = filepath.Join(xdg.CacheHome, This)
	Config = filepath.Join(xdg.ConfigHome, This)
)

func init() {
	err := os.MkdirAll(Cache, 0700)
	if err != nil {
		spool.Warn("%v\n", err)
	}
	err = os.MkdirAll(Config, 0700)
	if err != nil {
		spool.Warn("%v\n", err)
	}
}

func ReplaceHome(input string) string {
	return strings.Replace(input, Home, "~", 1)
}
