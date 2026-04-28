package shell

import (
	"errors"
	"pik/describe"
	"pik/model"
	"pik/runner"
	"pik/spool"
)

type Hydrated struct {
	runner.BaseHydration[*Target]
}

func (h *Hydrated) Icon() string {
	return "\uF489"
}

func (h *Hydrated) Description(src *model.HydratedSource) string {
	desc, err := describe.Describe(h.Target(), h.Target().File(src.Source))
	if err != nil {
		spool.Warn("%v\n", err)
		return ""
	}
	return desc
}

var WrongTargetError = errors.New("wrong target type")

func (s *shell) Hydrate(target model.Target) (model.HydratedTarget, error) {
	cast, ok := target.(*Target)
	if !ok {
		return nil, WrongTargetError
	}
	hyd := &Hydrated{BaseHydration: runner.Hydrated(cast)}
	return hyd, nil
}
