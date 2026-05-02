package shell

import (
	"github.com/ewy1/pik/describe"
	"github.com/ewy1/pik/model"
	"github.com/ewy1/pik/runner"
	"github.com/ewy1/pik/spool"
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
		_, _ = spool.Warn("%v\n", err)
		return ""
	}
	return desc
}
