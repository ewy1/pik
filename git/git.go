package git

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"pik/model"
	"pik/spool"
	"strconv"
	"strings"
)

type gitMod struct {
	Git string
	err error
}

func (g *gitMod) Init() error {
	p, err := exec.LookPath("git")
	if err != nil {
		g.err = err
		return nil
	}
	g.Git = p
	return nil
}

var Git = &gitMod{}

func (g *gitMod) Mod(source *model.Source, result *model.HydratedSource) error {
	gitFolder := filepath.Join(source.Path, ".git")
	if st, err := os.Stat(gitFolder); err == nil && st.IsDir() {
		if g.Git == "" {
			spool.Warn("source %v seems to be a git repository but git is not installed\n", source.Identity.Full)
			return nil
		}
		branch, err := g.Branch(source)
		if err != nil {
			spool.Warn("%v", err)
			return nil
		}
		ch, in, de, err := g.Diff(source)
		if err != nil {
			spool.Warn("%v", err)
			return nil
		}
		result.Git = &model.GitInfo{
			Branch:     branch,
			Insertions: in,
			Deletions:  de,
			Changes:    ch,
		}
	}

	return nil
}

func (g *gitMod) Branch(source *model.Source) (string, error) {
	cmd := exec.Command(g.Git, "branch", "--show-current")
	cmd.Dir = source.Path
	b, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(b)), err
}

var UnknownResponseError = errors.New("unknown response")

func (g *gitMod) Diff(source *model.Source) (int, int, int, error) {
	cmd := exec.Command(g.Git, "diff", "--shortstat")
	cmd.Dir = source.Path
	b, err := cmd.CombinedOutput()
	if err != nil {
		return 0, 0, 0, err
	}
	split := strings.Split(string(b), ",")
	changes := 0
	insertions := 0
	deletions := 0
	for _, s := range split {
		if strings.TrimSpace(s) == "" {
			return 0, 0, 0, nil
		}
		var e error
		pt := strings.Split(strings.TrimSpace(s), " ")
		num, e := strconv.Atoi(pt[0])
		switch {
		case strings.Contains(s, "changed"):
			changes = num
		case strings.Contains(s, "insertion"):
			insertions = num
		case strings.Contains(s, "deletion"):
			deletions = num
		default:
			return changes, insertions, deletions, UnknownResponseError
		}

		if e != nil {
			return changes, insertions, deletions, e
		}
	}
	return changes, insertions, deletions, nil
}
