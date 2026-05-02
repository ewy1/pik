package gnumake

import (
	"git.ewy.one/pik/model"
	"git.ewy.one/pik/runner"
	"os/exec"
)

type Target struct {
	runner.BaseTarget
	Name        string
	Description string
}

func (j *Target) File(src *model.Source) string {
	return Indexer.files[src.Path]
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

func (h *Hydrated) Description(src *model.HydratedSource) string {
	return h.Self.Description
}

func (h *Hydrated) Icon() string {
	return "\uE673"
}
