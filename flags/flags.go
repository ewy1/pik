package flags

import "github.com/spf13/pflag"

var (
	Here    = pflag.BoolP("here", "h", false, "run target in current directory instead of source location")
	At      = pflag.StringP("at", "@", "", "override run location")
	Single  = pflag.BoolP("single", "s", false, "do not run any triggers")
	All     = pflag.BoolP("all", "a", false, "get sources from cache instead of crawling")
	Dry     = pflag.BoolP("dry", "d", false, "print cmdlines instead of running them")
	Root    = pflag.BoolP("root", "r", false, "run targets (including triggers) with sudo")
	Yes     = pflag.BoolP("yes", "y", false, "auto-confirm y/n confirmations")
	Env     = pflag.StringArray("env", nil, "environment files or pre- or suffix")
	Version = pflag.BoolP("version", "v", false, "print version and exit")
	List    = pflag.BoolP("list", "l", false, "list available targets and exit")
	Inline  = pflag.BoolP("inline", "i", false, "do not use terminal alt screen")
)
