package shell

import (
	"os/exec"
	"path/filepath"
	"pik/model"
	"pik/runner"
)

type ShellTarget struct {
	runner.BaseTarget
	Shell    string
	Script   string
	SubValue []string
}

func (s *ShellTarget) String() string {
	return s.Label()
}

func (s *ShellTarget) Hydrate(_ *model.Source) (model.HydratedTarget, error) {
	return Runner.Hydrate(s)
}

func (s *ShellTarget) Sub() []string {
	return s.SubValue
}

func (s *ShellTarget) Label() string {
	return s.Identity.Full
}

func (s *ShellTarget) Create(src *model.Source) *exec.Cmd {
	return exec.Command(s.Shell, filepath.Join(src.Path, s.Script))
}
