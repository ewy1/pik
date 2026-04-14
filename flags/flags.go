package flags

import "github.com/spf13/pflag"

var (
	Here   = pflag.BoolP("here", "h", false, "run target in current directory instead of source location")
	At     = pflag.StringP("at", "@", "", "override run location")
	Single = pflag.BoolP("singgle", "s", false, "do not run any triggers")
	All    = pflag.BoolP("all", "a", false, "get sources from cache instead of crawling")
	Dry    = pflag.BoolP("dry", "d", false, "print cmdlines instead of running them")
	Root   = pflag.BoolP("root", "r", false, "run targets (including triggers) with sudo")
	Yes    = pflag.BoolP("yes", "y", false, "auto-confirm y/n confirmations")
	Env    = pflag.StringArray("env", nil, "environment files or pre- or suffix")
)
