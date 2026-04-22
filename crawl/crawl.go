package crawl

import (
	"path"
	"path/filepath"
	"slices"
	"strings"
)

func Evaluated(loc string) (string, error) {
	return filepath.EvalSymlinks(loc)
}

func RichLocations(origin string) []string {
	locs := Locations(origin)

	eval, err := Evaluated(origin)
	if err == nil && eval != origin {
		evaledLocations := Locations(eval)
		result := append(locs, evaledLocations...)
		result = slices.Compact(result)
		return result
	}
	return locs
}

func Locations(origin string) []string {
	origin = path.Clean(origin)
	var locs = []string{
		origin,
	}
	for {
		previous := locs[len(locs)-1]
		parent := ParentDir(previous)
		if previous == parent {
			break
		}
		locs = append(locs, parent)
	}
	return locs
}

func ParentDir(origin string) string {
	trimmedOrigin := strings.TrimSuffix(origin, "/")
	dir, _ := path.Split(trimmedOrigin)
	if dir == "" {
		return origin
	}
	return dir
}
