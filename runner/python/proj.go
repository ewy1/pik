package python

import (
	"github.com/ewy1/pik/model"
	"github.com/ewy1/pik/runner"
	"os/exec"
	"path/filepath"
)

type Project struct {
	runner.BaseTarget
	Cmd string
}

func (p *Project) File(src *model.Source) string {
	return Python.files[src.Path]
}

type Hydrated struct {
	runner.BaseHydration[*Project]
}

func (h *Hydrated) Description(src *model.HydratedSource) string {
	return h.Self.Cmd
}

func (h *Hydrated) Icon() string {
	return "\uE606"
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
