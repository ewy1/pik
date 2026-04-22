package search

import (
	"pik/model"
	"slices"
)

// Result is a struct containing information about the search and search results
type Result struct {
	// Target is the target selected by the search
	Target model.Target
	// Source is the source belonging to the selected Target
	Source *model.Source
	// NeedsConfirmation is true when there are discrepancies between expected and actual invocation
	NeedsConfirmation bool
	// Overridden is whether it was overridden by a .override target
	Overridden bool
	// Sub is the subcategory or -folder
	Sub []string
	// Args are the remaining arguments which we should pass to the target
	Args []string
}

// Search is the meat of pik
// since there are a ton of different ways to invoke targets, leave a unit test
// when you change this
func Search(s *model.State, args ...string) *Result {
	var target model.Target
	var targetSource *model.Source
	var confirm bool
	var overridden bool
	var subdir []string
	var forward []string
	var suspect model.Target
	var suspectSource *model.Source

args_loop:
	for i, arg := range args {
		for _, src := range s.Sources {

			if targetSource == nil {
				if src.Is(arg) {
					targetSource = src

					// only try to find the default target if this is the last argument
					if len(args)-1 != i {
						continue args_loop
					}

					// try to look for arg target with the same name as the source
					// "default target" of sorts
					for _, t := range targetSource.Targets {
						if t.Matches(arg) {
							target = t
							continue args_loop
						}
					}

					continue args_loop
				}
			}

			if target == nil && targetSource == nil {

				// uncertain about source, check ours to see if any match
				for _, t := range src.Targets {
					if t.Matches(arg) {
						if slices.Equal(t.Sub(), subdir) {
							target = t
							targetSource = src
						} else {
							suspect = t
							suspectSource = src
						}
						continue args_loop
					}
				}

			} else if target == nil { // && targetSource == nil (but it is always true)

				// source located,
				for _, t := range targetSource.Targets {
					if t.Matches(arg) {
						target = t
						continue args_loop
					}
				}
				// if we find the right target
				for _, t := range src.Targets {
					if t.Matches(arg) {
						confirm = true
						suspect = t
						suspectSource = src
						continue args_loop
					}
				}
			}

		}

		if target == nil && suspect == nil {
			subdir = append(subdir, arg)
			continue args_loop
		} else if targetSource != nil || suspect != nil {
			forward = append(forward, arg)
			continue args_loop
		}
	}

	if suspect != nil && target == nil {
		target = suspect
		targetSource = suspectSource

		if !(suspect.Sub() != nil && subdir == nil) {
			confirm = true
		}
	}

	if target != nil && target.Sub() != nil && subdir != nil && !slices.Equal(target.Sub(), subdir) {
		confirm = true
	}

	if target == nil {
		forward = args
	}

	if target != nil && targetSource != nil {
		for _, t := range targetSource.Targets {
			if slices.Equal(t.Invocation(targetSource), target.Invocation(targetSource)) {
				if t.Tags().Has(model.Override) {
					overridden = true
					target = t
				}
			}
		}
	}

	return &Result{
		Target:            target,
		Source:            targetSource,
		NeedsConfirmation: confirm,
		Overridden:        overridden,
		Sub:               subdir,
		Args:              forward,
	}
}
