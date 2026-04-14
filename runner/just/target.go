package just

import (
	"os/exec"
	"pik/model"
	"pik/runner"
)

type Target struct {
	runner.BaseTarget
	Category string
}

func (j Target) Create(s *model.Source) *exec.Cmd {
	return exec.Command(Indexer.path, j.Identity.Full)
}

func (j Target) Sub() []string {
	if j.Category != "" {
		return []string{j.Category}
	}
	return nil
}

func (j Target) Label() string {
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
	return "\uF039"
}
