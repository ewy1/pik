package run

import (
	"fmt"
	"os"
	"pik/flags"
	"pik/menu"
	"pik/model"
	"slices"
)

func Run(source *model.Source, target model.Target, args ...string) error {
	tags := target.Tags()
	skipTriggers := tags.Has(model.Single) || *flags.Single

	if !skipTriggers {
		err := Pre(source, target)
		if err != nil {
			return err
		}
	}
	err := Exec(source, target, args...)
	fmt.Println()
	if err != nil {
		return err
	}
	if !skipTriggers {
		err := Post(source, target)
		if err != nil {
			return err
		}
		err = Final(source, target)
		if err != nil {
			return err
		}
	}
	return nil
}

func Pre(source *model.Source, target model.Target) error {
	return ExecWithTrigger(source, target, model.Pre)
}

func Post(source *model.Source, target model.Target) error {
	return ExecWithTrigger(source, target, model.Post)
}

func Final(source *model.Source, target model.Target) error {
	return ExecWithTrigger(source, target, model.Final)
}

func ExecWithTrigger(source *model.Source, target model.Target, tag model.Tag) error {
	for _, t := range source.Targets {
		if t.Tags().Has(tag) {
			triggerSub := t.Sub()
			targetSub := target.Sub()

			for _, targetSubPart := range triggerSub {
				if !slices.Contains(targetSub, targetSubPart) {
					continue
				}
			}

			err := Exec(source, t)
			fmt.Println()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Exec(source *model.Source, target model.Target, args ...string) error {
	_, _ = fmt.Fprint(os.Stderr, menu.Banner(source, target, args...))
	loc := source.Path
	tags := target.Tags()
	if *flags.At != "" {
		loc = *flags.At
	} else if tags.Has(model.Here) || *flags.Here {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		loc = wd
	}
	cmd := target.Create(source)
	cmd.Dir = loc
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Args = append(cmd.Args, args...)

	if *flags.Dry {
		_, _ = fmt.Fprintln(os.Stderr, menu.InlineCmd(cmd))
		return nil
	}

	if *flags.Root {
		cmd.Args = append([]string{"sudo"}, cmd.Args...)
	}
	return cmd.Run()
}
