package shell

import (
	"github.com/ewy1/pik/model"
	"github.com/ewy1/pik/runner"
	"os/exec"
)

type Target struct {
	runner.BaseTarget
	Shell  string
	Script string
}

func (s *Target) File(src *model.Source) string {
	return s.Script
}

func (s *Target) String() string {
	return s.Label()
}

func (s *Target) Hydrate(_ *model.Source) (model.HydratedTarget, error) {
	return Runner.Hydrate(s)
}

func (s *Target) Sub() []string {
	return s.BaseTarget.MySub
}

func (s *Target) Label() string {
	return s.Identity.Full
}

func (s *Target) Create(_ *model.Source) *exec.Cmd {
	return exec.Command(s.Shell, s.Script)
}
