package env

import (
	"github.com/joho/godotenv"
	"io/fs"
	"os"
	"path/filepath"
	"pik/flags"
	"pik/indexers/pikdex"
	"pik/model"
	"pik/spool"
	"slices"
)

// IsEnv returns whether a given file is suitable for environment loading
// this method respects the --env flag
func IsEnv(file string) bool {
	options := []string{
		".env",
	}
	for _, e := range *flags.Env {
		options = append(options,
			".env-"+e,
			".env."+e,
			e+".env",
			"."+e+".env")
	}
	return slices.Contains(options, file)
}

// Files returns a list of files (that exist) that should be indexed and used as environment files
func Files(f fs.FS, p string, deep bool) []string {
	var result []string
	dir, err := fs.ReadDir(f, p)
	if err != nil {
		return nil
	}
	for _, e := range dir {
		if e.IsDir() && slices.Contains(pikdex.Roots, e.Name()) && deep {
			result = append(result, Files(f, e.Name(), false)...)
		}
		if !e.IsDir() && IsEnv(e.Name()) {
			result = append(result, filepath.Join(p, e.Name()))
		}
	}
	return result
}

// Get returns all environment key-value pairs we should index for a source
func Get(src *model.Source) []string {
	f := os.DirFS(src.Path)
	var result []string
	files := Files(f, ".", true)
	for _, f := range files {
		res, err := godotenv.Read(filepath.Join(src.Path, f))
		if err != nil {
			spool.Warn("%v", err)
			continue
		}
		for k, v := range res {
			result = append(result, k+"="+v)
		}
	}
	return result
}
