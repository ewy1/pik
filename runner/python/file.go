package python

import (
	"git.ewy.one/pik.git/describe"
	"git.ewy.one/pik.git/model"
	"git.ewy.one/pik.git/runner"
	"git.ewy.one/pik.git/spool"
	"os/exec"
	"path/filepath"
)

type File struct {
	runner.BaseTarget
	Path string
}

func (p *File) File(src *model.Source) string {
	return p.Path
}

type HydratedFileTarget struct {
	runner.BaseHydration[*File]
}

func (h *HydratedFileTarget) Description(src *model.HydratedSource) string {
	desc, err := describe.Describe(h.Target(), h.Self.Path)
	if err != nil {
		spool.Warn("%v\n", err)
	}
	return desc
}

func (h *HydratedFileTarget) Icon() string {
	return "\uE606"
}

func (p *File) Create(s *model.Source) *exec.Cmd {
	var cmd []string
	if Python.Uv != "" {
		cmd = []string{Python.Uv, "run", "--", p.Path}
	} else if venv := Python.VenvFor(s); venv != "" {
		cmd = []string{filepath.Join(s.Path, venv, "bin", "python3"), p.Path}
	} else {
		sysPath, err := exec.LookPath("python3")
		if err != nil {
			return nil
		}
		cmd = []string{sysPath, p.Path}
	}
	return exec.Command(cmd[0], cmd[1:]...)
}

func (p *File) Sub() []string {
	return nil
}

func (p *File) Label() string {
	return p.Full
}

func (p *File) Hydrate(src *model.Source) (model.HydratedTarget, error) {
	return &HydratedFileTarget{
		BaseHydration: runner.Hydrated(p),
	}, nil
}
