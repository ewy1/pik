package model

import (
	"pik/identity"
	"pik/paths"
	"pik/spool"
)

// Source is a location containing stuff we can run
// these get created when we find a makefile, .pik folder, etc.
type Source struct {
	identity.Identity
	Tags
	Path    string
	Targets []Target
}

// HydratedSource is a Source with additional hydration
// for the menu.
// these do not get created unless we show the menu
type HydratedSource struct {
	*Source
	HydratedTargets []HydratedTarget
	Aliases         []string
	Icon            string
	Git             *GitInfo
}

func (s *Source) Label() string {
	return s.Identity.Full
}

func (s *HydratedSource) Label() string {
	if len(s.Aliases) > 0 {
		return s.Aliases[0]
	}
	return s.Identity.Full
}

func (s *Source) Hydrate(hydrators []Modder) *HydratedSource {
	hs := &HydratedSource{
		Source:          s,
		HydratedTargets: make([]HydratedTarget, 0, len(s.Targets)),
	}
	for _, h := range hydrators {
		err := h.Mod(s, hs)
		if err != nil {
			spool.Warn("%v", err)
		}
	}
	for _, t := range s.Targets {
		if !t.Visible() {
			continue
		}
		ht, err := t.Hydrate(s)
		if err != nil {
			spool.Warn("%v", err)
			continue
		}
		hs.HydratedTargets = append(hs.HydratedTargets, ht)
	}
	return hs
}

func (s *Source) ShortPath() string {
	return paths.ReplaceHome(s.Path)
}
