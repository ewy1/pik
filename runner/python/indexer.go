package python

import (
	"git.ewy.one/pik.git/model"
	"github.com/pelletier/go-toml/v2"
	"io/fs"
	"os"
	"path/filepath"
)

type pyproj struct {
	Project struct {
		Scripts map[string]string
	}
}

func (p *python) Index(path string, f fs.FS, runners []model.Runner) ([]model.Target, error) {
	for _, pt := range VenvPaths {
		if stat, err := fs.Stat(f, filepath.Join(pt)); err == nil {
			if stat.IsDir() {
				p.Venvs[path] = filepath.Join(path, pt)
			}
		}
	}
	content, err := fs.ReadFile(f, "pyproject.toml")
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	p.files[path] = filepath.Join(path, "pyproject.toml")

	pp := &pyproj{}

	err = toml.Unmarshal(content, pp)
	if err != nil {
		return nil, err
	}

	var targets = make([]model.Target, 0, len(pp.Project.Scripts))
	for n, s := range pp.Project.Scripts {
		targets = append(targets, Python.CreateProjTarget(n, s))
	}
	return targets, nil
}
