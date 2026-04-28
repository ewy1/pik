package shell

import (
	"os/exec"
	"pik/model"
	"pik/runner"
)

type Target struct {
	runner.BaseTarget
	Shell  string
	Script string
}

func (s *Target) String() string {
	return s.Label()
}

func (s *Target) Hydrate(_ *model.Source) (model.HydratedTarget, error) {
	return Runner.Hydrate(s)
}

func (s *Target) Sub() []string {
	return s.BaseTarget.Sub
}

func (s *Target) Label() string {
	return s.Identity.Full
}

func (s *Target) Create(_ *model.Source) *exec.Cmd {
	return exec.Command(s.Shell, s.Script)
}
