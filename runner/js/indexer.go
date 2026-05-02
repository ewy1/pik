package js

import (
	"encoding/json"
	"errors"
	"github.com/ewy1/pik/model"
	"io/fs"
	"os"
	"path/filepath"
)

type Package struct {
	Scripts map[string]string `json:"scripts"`
}

func (n *js) Index(path string, f fs.FS, runners []model.Runner) ([]model.Target, error) {
	p := &Package{}
	// are there any other package.jsons? i hope not, because i don't know them
	content, err := os.ReadFile(filepath.Join(path, "package.json"))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, p)
	if err != nil {
		return nil, err
	}

	var targets []model.Target
	if n.Npm == "" {
		return nil, NoNpm
	}
	for k, s := range p.Scripts {
		targets = append(targets, n.CreateRun(k, s))
	}

	return targets, nil
}
