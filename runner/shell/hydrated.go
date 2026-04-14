package shell

import (
	"pik/runner"
)

type HydratedShellTarget struct {
	runner.BaseHydration[*ShellTarget]
}

func (h *HydratedShellTarget) Icon() string {
	return "\uF489"
}

func (h *HydratedShellTarget) Description() string {
	return "//TODO"
}
