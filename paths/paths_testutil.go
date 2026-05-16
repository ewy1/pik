//go:build test

package paths

var backups = make(map[*Path]string)

func SetAll(folder string) {
	for _, p := range Paths {
		Set(p, folder)
	}
}

// Set temporarily sets the paths for unit test purposes
// remember to defer Reset
func Set(target *Path, value string) {
	backups[target] = target.String()
	target.Set(value)
}

// Reset sets the path variables back to before the unit test
func Reset() {
	for path, oldValue := range backups {
		path.Set(oldValue)
	}
	backups = make(map[*Path]string)
}
