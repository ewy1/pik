package gnumake

import (
	"errors"
	"io/fs"
	"os/exec"
	"pik/identity"
	"pik/model"
	"pik/runner"
	"regexp"
	"slices"
	"strings"
)

type make struct {
	path string
}

var Indexer = &make{}

var Makefiles = []string{
	"Makefile",
	"makefile",
}

func (m *make) Index(path string, f fs.FS, _ []model.Runner) ([]model.Target, error) {

	entries, err := fs.ReadDir(f, ".")
	if err != nil {
		return nil, err
	}
	makefile := ""
	for _, e := range entries {
		if !e.IsDir() && slices.Contains(Makefiles, strings.ToLower(e.Name())) {
			content, err := fs.ReadFile(f, e.Name())
			if err != nil {
				return nil, err
			}
			makefile = string(content)
			break
		}
	}

	if makefile == "" {
		return nil, nil
	}

	err = m.findMake()
	if err != nil {
		return nil, err
	}

	return ParseOutput(makefile), nil
}

var makeRegex = regexp.MustCompile("^([a-zA-Z-]*):((.*?)# (.*))?")

func ParseOutput(input string) []model.Target {
	var targets []string
	match := makeRegex.FindAllString(input, len(input))
	for _, m := range match {
		targets = append(targets, m)
	}

	var result []model.Target
	for _, t := range targets {
		split := strings.SplitN(t, "#", 2)
		name := split[0]
		name = strings.TrimSpace(name)
		name = strings.TrimSuffix(name, ":")
		tgt := &Target{
			BaseTarget: runner.BaseTarget{
				Identity: identity.New(name),
				MyTags:   nil,
			},
			Name: name,
		}
		if len(split) > 1 {
			tgt.Description = strings.TrimSpace(split[1])
		}
		result = append(result, tgt)
	}
	return result
}

var NoJustError = errors.New("no make in $PATH but source contains makefile")

func (m *make) findMake() error {
	loc, err := exec.LookPath("make")
	if errors.Is(err, exec.ErrNotFound) {
		return NoJustError
	} else if err != nil {
		return err
	}
	m.path = loc
	return nil
}
