package python

import (
	"os/exec"
	"path/filepath"
	"pik/model"
	"pik/runner"
)

type FileTarget struct {
	runner.BaseTarget
	File string
}

type HydratedFileTarget struct {
	runner.BaseHydration[*FileTarget]
}

func (h *HydratedFileTarget) Icon() string {
	return "\uE606"
}

func (p *FileTarget) Create(s *model.Source) *exec.Cmd {
	var cmd []string
	if Python.Uv != "" {
		cmd = []string{Python.Uv, "run", "--", p.File}
	} else if venv := Python.VenvFor(s); venv != "" {
		cmd = []string{filepath.Join(s.Path, venv, "bin", "python3"), p.File}
	} else {
		sysPath, err := exec.LookPath("python3")
		if err != nil {
			return nil
		}
		cmd = []string{sysPath, p.File}
	}
	return exec.Command(cmd[0], cmd[1:]...)
}

func (p *FileTarget) Sub() []string {
	return nil
}

func (p *FileTarget) Label() string {
	return p.Full
}

func (p *FileTarget) Hydrate(src *model.Source) (model.HydratedTarget, error) {
	return &HydratedFileTarget{
		BaseHydration: runner.Hydrated(p),
	}, nil
}
