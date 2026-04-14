package shell

import (
	"pik/describe"
	"pik/runner"
	"pik/spool"
)

type Hydrated struct {
	runner.BaseHydration[*Target]
}

func (h *Hydrated) Icon() string {
	return "\uF489"
}

func (h *Hydrated) Description() string {
	desc, err := describe.Describe(h.BaseTarget.Script)
	if err != nil {
		spool.Warn("%v\n", err)
		return ""
	}
	return desc
}
