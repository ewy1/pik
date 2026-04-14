package gnumake

import (
	"os/exec"
	"pik/model"
	"pik/runner"
)

type Target struct {
	runner.BaseTarget
	Name        string
	Description string
}

func (j *Target) Create(s *model.Source) *exec.Cmd {
	return exec.Command(Indexer.path, j.Identity.Full)
}

var makeSub = []string{
	"make",
}

func (j *Target) Sub() []string {
	return makeSub
}

func (j *Target) Label() string {
	return j.Identity.Full
}

func (j *Target) Hydrate(src *model.Source) (model.HydratedTarget, error) {
	return &Hydrated{
		BaseHydration: runner.Hydrated(j),
	}, nil
}

type Hydrated struct {
	runner.BaseHydration[*Target]
}

func (h *Hydrated) Icon() string {
	return "\uE673"
}

func (h *Hydrated) Description() string {
	return h.BaseTarget.Description
}
