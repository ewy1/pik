package just

import (
	"os/exec"
	"pik/model"
	"pik/runner"
)

type JustTarget struct {
	runner.BaseTarget
	Category string
}

func (j JustTarget) Create(s *model.Source) *exec.Cmd {
	return exec.Command(Indexer.path, j.Identity.Full)
}

func (j JustTarget) Sub() []string {
	if j.Category != "" {
		return []string{j.Category}
	}
	return nil
}

func (j JustTarget) Label() string {
	return j.Identity.Full
}

func (j *JustTarget) Hydrate(src *model.Source) (model.HydratedTarget, error) {
	return &HydratedJustTarget{
		BaseHydration: runner.Hydrated(j),
	}, nil
}

type HydratedJustTarget struct {
	runner.BaseHydration[*JustTarget]
}

func (h *HydratedJustTarget) Icon() string {
	return "\uF039"
}
