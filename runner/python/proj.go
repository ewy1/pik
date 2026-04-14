package python

import (
	"os/exec"
	"path/filepath"
	"pik/model"
	"pik/runner"
)

type Project struct {
	runner.BaseTarget
	Cmd string
}

type Hydrated struct {
	runner.BaseHydration[*Project]
}

func (h *Hydrated) Icon() string {
	return "\uE606"
}

func (h *Hydrated) Description() string {
	return h.BaseTarget.Cmd
}

func (p *Project) Create(s *model.Source) *exec.Cmd {
	var cmd []string
	if Python.Uv != "" {
		cmd = []string{Python.Uv, "run", "--", p.Cmd}
	} else if venv := Python.VenvFor(s); venv != "" {
		cmd = []string{filepath.Join(s.Path, venv, "bin", "python"), p.Cmd}
	}
	return exec.Command(cmd[0], cmd[1:]...)
}

func (p *Project) Sub() []string {
	return nil
}

func (p *Project) Label() string {
	return p.Cmd
}

func (p *Project) Hydrate(src *model.Source) (model.HydratedTarget, error) {
	return &Hydrated{
		BaseHydration: runner.Hydrated(p),
	}, nil
}
