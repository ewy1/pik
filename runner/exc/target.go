package exc

import (
	"os/exec"
	"pik/describe"
	"pik/model"
	"pik/runner"
	"pik/spool"
)

type Executable struct {
	runner.BaseTarget
	Path string
}

type Hydrated struct {
	*runner.BaseHydration[*Executable]
}

func (h *Hydrated) Icon() string {
	return "\uEAE8"
}

func (h *Hydrated) Description(src *model.HydratedSource) string {
	d, err := describe.Describe(h.Self, h.Self.Path)
	if err != nil {
		spool.Warn("%v\n", err)
	}
	return d
}

func (e *Executable) Create(s *model.Source) *exec.Cmd {
	return exec.Command(e.Path)
}

func (e *Executable) Label() string {
	return e.Identity.Full
}

func (e *Executable) Hydrate(src *model.Source) (model.HydratedTarget, error) {
	return &Hydrated{
		BaseHydration: &runner.BaseHydration[*Executable]{
			Self: e,
		},
	}, nil
}

func (e *Executable) File(src *model.Source) string {
	return e.Path
}
