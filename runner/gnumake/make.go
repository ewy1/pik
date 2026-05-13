package gnumake

import (
	"errors"
	"github.com/ewy1/pik/identity"
	"github.com/ewy1/pik/model"
	"github.com/ewy1/pik/runner"
	"io/fs"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

type gnumake struct {
	path  string
	files map[string]string
}

var Indexer = &gnumake{
	files: make(map[string]string),
}

var Makefiles = []string{
	"Makefile",
	"makefile",
}

func (m *gnumake) Index(path string, f fs.FS, _ []model.Runner) ([]model.Target, error) {

	entries, err := fs.ReadDir(f, ".")
	if err != nil {
		return nil, err
	}
	makefile := ""
	for _, e := range entries {
		if !e.IsDir() && slices.Contains(Makefiles, strings.ToLower(e.Name())) {
			m.files[path] = filepath.Join(path, e.Name())
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

var NoMakeError = errors.New("source contains makefile but no gnumake in $PATH")

func (m *gnumake) findMake() error {
	loc, err := exec.LookPath("make")
	if errors.Is(err, exec.ErrNotFound) {
		return NoMakeError
	} else if err != nil {
		return err
	}
	m.path = loc
	return nil
}
