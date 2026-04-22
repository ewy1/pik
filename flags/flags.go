package flags

import "github.com/spf13/pflag"

var (
	// Here makes pik run the target at the current location instead of in its source directory
	Here = pflag.BoolP("here", "h", false, "run target in current directory instead of source location")
	// At makes pik run in an arbitrary location
	At = pflag.StringP("at", "@", "", "override run location")
	// Single will make pik skip triggers
	Single = pflag.BoolP("single", "s", false, "do not run any triggers")
	// All will make pik load cached state instead of crawling
	All = pflag.BoolP("all", "a", false, "get sources from cache instead of crawling")
	// Dry will prevent pik from running targets and output the command as text instead
	Dry = pflag.BoolP("dry", "d", false, "print cmdlines instead of running them")
	// Root will prefix the target with sudo to run it as root
	Root = pflag.BoolP("root", "r", false, "run targets (including triggers) with sudo")
	// Yes will skip y/n prompts and always answer yes
	Yes = pflag.BoolP("yes", "y", false, "auto-confirm y/n confirmations")
	// Env can be used to load additional .env files
	Env = pflag.StringArray("env", nil, "environment files or pre- or suffix")
	// Version means we should print the pik version and exit
	Version = pflag.BoolP("version", "v", false, "print version and exit")
	// List means we should output available targets separated by $IFS
	List = pflag.BoolP("list", "l", false, "list available targets and exit")
	// Inline means pik does not go to the terminal alt screen
	Inline = pflag.BoolP("inline", "i", false, "do not use terminal alt screen")
)
