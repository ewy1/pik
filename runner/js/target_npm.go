package js

import (
	"github.com/ewy1/pik/model"
	"github.com/ewy1/pik/runner"
	"os/exec"
	"path/filepath"
)

type Npm struct {
	runner.BaseTarget
	Name string
	Cmd  string
}

func (n *Npm) Icon() string {
	return "\uE60B"
}

func (n *Npm) Description(src *model.HydratedSource) string {
	return n.Cmd
}

func (n *Npm) Target() model.Target {
	return n
}

func (n *Npm) Create(s *model.Source) *exec.Cmd {
	return exec.Command(Js.Npm, "run", n.Name)
}

func (n *Npm) Label() string {
	return n.Name
}

func (n *Npm) Hydrate(src *model.Source) (model.HydratedTarget, error) {
	return n, nil
}

func (n *Npm) File(src *model.Source) string {
	return filepath.Join(src.Path, "package.json")
}
