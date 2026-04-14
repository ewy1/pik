package python

import (
	"errors"
	"io/fs"
	"os/exec"
	"path/filepath"
	"pik/identity"
	"pik/model"
	"pik/runner"
)

type python struct {
	Venvs  map[string]string
	Uv     string
	System string
}

func (p python) Init() error {
	uv, err := exec.LookPath("uv")
	if err != nil && !errors.Is(err, exec.ErrNotFound) {
		return err
	}
	p.Uv = uv
	sys, err := exec.LookPath("python3")
	if err == nil {
		p.System = sys
	}
	return err
}

func (p python) Hydrate(target model.Target) (model.HydratedTarget, error) {
	//TODO implement me
	panic("implement me")
}

func (p python) Wants(fs fs.FS, file string, entry fs.DirEntry) (bool, error) {
	return !entry.IsDir() && filepath.Ext(entry.Name()) == ".py", nil
}

func (p python) VenvFor(src *model.Source) string {
	venvPath := p.Venvs[src.Path]
	if venvPath != "" {
		return venvPath
	}
	return ""
}

func (p python) PyFor(src *model.Source) []string {
	if p.Uv != "" {
		return []string{p.Uv, "run", "--"}
	}
	if venv := p.VenvFor(src); venv != "" {
		return []string{filepath.Join(src.Path, venv, "bin", "python")}
	}
	return nil
}

func (p python) CreateProjTarget(name string, cmd string) model.Target {
	return &Project{
		BaseTarget: runner.BaseTarget{
			Identity: identity.New(name),
		},
		Cmd: cmd,
	}
}

func (p python) CreateTarget(fs fs.FS, source string, file string, entry fs.DirEntry) (model.Target, error) {
	_, filename := filepath.Split(file)
	return &File{
		BaseTarget: runner.BaseTarget{
			Identity: identity.New(filename),
			MyTags:   model.TagsFromFilename(filename),
		},
		File: file,
	}, nil
}

var VenvPaths = []string{
	".venv",
	"venv",
}

var Python = &python{
	Venvs: map[string]string{},
}
