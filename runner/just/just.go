package just

import (
	"errors"
	"io/fs"
	"os/exec"
	"pik/identity"
	"pik/model"
	"pik/runner"
	"strings"
)

type just struct {
	path string
}

var Indexer = &just{}

func (j *just) Index(path string, f fs.FS, runners []model.Runner) ([]model.Target, error) {

	entries, err := fs.ReadDir(f, ".")
	if err != nil {
		return nil, err
	}
	hasJustfile := false
	for _, e := range entries {
		if !e.IsDir() && strings.ToLower(e.Name()) == "justfile" {
			hasJustfile = true
			break
		}
	}

	if !hasJustfile {
		return nil, nil
	}

	err = j.findJust()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(j.path, "--list")
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return ParseOutput(string(out)), nil
}

func ParseOutput(input string) []model.Target {
	categories := make(map[string][]string)
	currentCategory := ""
	for _, line := range strings.Split(input, "\n") {
		// strip comment
		line = strings.SplitN(line, "#", 2)[0]
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentCategory = line[1 : len(line)-1]
			continue
		}

		tgt := strings.SplitN(line, " ", 2)[0]

		if tgt == "" {
			continue
		}

		categories[currentCategory] = append(categories[currentCategory], tgt)
	}

	var result []model.Target
	for c, targets := range categories {
		for _, t := range targets {
			result = append(result, &JustTarget{
				BaseTarget: runner.BaseTarget{
					Identity: identity.New(t),
				},
				Category: c,
			})
		}
	}
	return result
}

var NoJustError = errors.New("no just in $PATH but source contains justfile")

func (j *just) findJust() error {
	loc, err := exec.LookPath("just")
	if errors.Is(err, exec.ErrNotFound) {
		return NoJustError
	} else if err != nil {
		return err
	}
	j.path = loc
	return nil
}
