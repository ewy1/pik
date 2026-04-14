package gnumake

import (
	"os/exec"
	"pik/model"
	"pik/runner"
)

type MakeTarget struct {
	runner.BaseTarget
	Name        string
	Description string
}

func (j *MakeTarget) Create(s *model.Source) *exec.Cmd {
	return exec.Command(Indexer.path, j.Identity.Full)
}

var makeSub = []string{
	"make",
}

func (j *MakeTarget) Sub() []string {
	return makeSub
}

func (j *MakeTarget) Label() string {
	return j.Identity.Full
}

func (j *MakeTarget) Hydrate(src *model.Source) (model.HydratedTarget, error) {
	return &HydratedJustTarget{
		BaseHydration: runner.Hydrated(j),
	}, nil
}

type HydratedJustTarget struct {
	runner.BaseHydration[*MakeTarget]
}

func (h *HydratedJustTarget) Icon() string {
	return "\uE673"
}
