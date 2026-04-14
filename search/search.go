package search

import (
	"pik/model"
	"slices"
)

func Search(s *model.State, args ...string) (model.Target, *model.Source, bool, []string, []string) {
	var target model.Target
	var suspect model.Target
	var suspectSource *model.Source
	var targetSource *model.Source
	var forward []string
	var subdir []string
	confirm := false

args_loop:
	for _, a := range args {
		for _, src := range s.Sources {

			if targetSource == nil {
				if src.Is(a) {
					targetSource = src
					for _, t := range targetSource.Targets {
						if t.Matches(a) {
							target = t
							continue args_loop
						}
					}
					continue args_loop
				}
			}

			if target == nil && targetSource == nil {
				for _, t := range src.Targets {
					if t.Matches(a) {
						target = t
						targetSource = src
						continue args_loop
					}
				}
			} else if target == nil { // && targetSource == nil (but it is always true)
				for _, t := range targetSource.Targets {
					if t.Matches(a) {
						target = t
						continue args_loop
					}
				}
				// if we find the right target
				for _, t := range src.Targets {
					if t.Matches(a) {
						confirm = true
						suspect = t
						suspectSource = src
						continue args_loop
					}
				}
			}

		}

		if target == nil {
			subdir = append(subdir, a)
			continue args_loop
		} else if targetSource != nil {
			forward = append(forward, a)
			continue args_loop
		}
	}

	if suspect != nil && target == nil {
		target = suspect
		targetSource = suspectSource
		confirm = true
	}

	if target != nil && target.Sub() != nil && subdir != nil && !slices.Equal(target.Sub(), subdir) {
		confirm = true
	}

	if target == nil {
		forward = args
	}

	return target, targetSource, confirm, subdir, forward
}
