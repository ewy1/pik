package shell

import (
	"pik/describe"
	"pik/runner"
	"pik/spool"
)

type HydratedShellTarget struct {
	runner.BaseHydration[*ShellTarget]
}

func (h *HydratedShellTarget) Icon() string {
	return "\uF489"
}

func (h *HydratedShellTarget) Description() string {
	desc, err := describe.Describe(h.BaseTarget.Script)
	if err != nil {
		spool.Warn("%v\n", err)
		return ""
	}
	return desc
}
