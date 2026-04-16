package pikdex

import (
	"pik/model"
	"strings"
)

func (u *pikdex) Mod(src *model.Source, result *model.HydratedSource) error {
	mod := u.mods[strings.TrimSuffix(src.Path, "/")]
	if mod.Path != "" {
		if mod.Aliases != nil {
			result.Aliases = append(result.Aliases, mod.Aliases...)
		}
		if mod.Icon != "" {
			result.Icon = mod.Icon
		}
	}
	return nil
}
