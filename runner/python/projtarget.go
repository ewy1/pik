package python

import (
	"os/exec"
	"path/filepath"
	"pik/model"
	"pik/runner"
)

type ProjTarget struct {
	runner.BaseTarget
	Cmd string
}

type HydratedProjTarget struct {
	runner.BaseHydration[*ProjTarget]
}

func (h *HydratedProjTarget) Icon() string {
	return "\uE606"
}

func (h *HydratedProjTarget) Description() string {
	return h.BaseTarget.Cmd
}

func (p *ProjTarget) Create(s *model.Source) *exec.Cmd {
	var cmd []string
	if Python.Uv != "" {
		cmd = []string{Python.Uv, "run", "--", p.Cmd}
	} else if venv := Python.VenvFor(s); venv != "" {
		cmd = []string{filepath.Join(s.Path, venv, "bin", "python"), p.Cmd}
	}
	return exec.Command(cmd[0], cmd[1:]...)
}

func (p *ProjTarget) Sub() []string {
	return nil
}

func (p *ProjTarget) Label() string {
	return p.Cmd
}

func (p *ProjTarget) Hydrate(src *model.Source) (model.HydratedTarget, error) {
	return &HydratedProjTarget{
		BaseHydration: runner.Hydrated(p),
	}, nil
}
