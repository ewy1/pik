package model

import (
	"pik/identity"
	"pik/paths"
	"pik/spool"
)

type Source struct {
	identity.Identity
	Tags
	Path    string
	Targets []Target
}

type HydratedSource struct {
	*Source
	HydratedTargets []HydratedTarget
	Aliases         []string
	Icon            string
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

func (s *Source) Hydrate(hydrators []Hydrator) *HydratedSource {
	hs := &HydratedSource{
		Source:          s,
		HydratedTargets: make([]HydratedTarget, 0, len(s.Targets)),
	}
	for _, h := range hydrators {
		err := h.Hydrate(s, hs)
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
