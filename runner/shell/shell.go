package shell

import (
	"bufio"
	"errors"
	"io/fs"
	"os/exec"
	"path/filepath"
	"pik/identity"
	"pik/indexers/pikdex"
	"pik/model"
	"pik/runner"
	"pik/spool"
	"slices"
	"strings"
)

//TODO: Clean up shell selection? Maybe default to bash?

var NoContentError = errors.New("not enough content in target")
var NoShellError = errors.New("could not find any shell interpreters")

var ExtShellMap = map[string]string{
	".sh":  "bash",
	".ps1": "powershell",
}

var Shells = []string{"bash", "bash.exe", "zsh", "fish", "powershell", "powershell.exe", "cmd.exe"}

var Runner = &shell{
	Locations: map[string]string{},
}

type shell struct {
	Locations map[string]string
}

func (s *shell) Wants(f fs.FS, file string, entry fs.DirEntry) (bool, error) {
	if entry != nil && entry.IsDir() {
		return false, nil
	}

	fd, err := f.Open(file)
	if err != nil {
		return false, err
	}
	scanner := bufio.NewScanner(fd)
	scanner.Split(bufio.ScanRunes)
	if !scanner.Scan() {
		return false, nil
	}
	txt := scanner.Text()
	if txt == "#" { //
		return true, nil
	}
	for k, _ := range ExtShellMap {
		if strings.HasSuffix(file, k) {
			return true, nil
		}
	}
	return false, nil
}

func (s *shell) Find(shell string) (string, error) {
	if s.Locations[shell] != "" {
		return s.Locations[shell], nil
	}

	if p, err := exec.LookPath(shell); err == nil {
		s.Locations[shell] = p
		return shell, nil
	} else {
		return "", err
	}
}

func (s *shell) CreateTarget(fs fs.FS, src string, file string, _ fs.DirEntry) (model.Target, error) {
	shell, err := s.ShellFor(fs, file)
	if err != nil {
		return nil, err
	}
	_, filename := filepath.Split(file)
	var sub []string
	split := strings.Split(file, "/")
	for _, p := range split {
		if slices.Contains(pikdex.Roots, p) {
			continue
		}
		if filename == p {
			continue
		}
		sub = append(sub, p)
	}
	return &Target{
		BaseTarget: runner.BaseTarget{
			Identity: identity.New(filename),
			MyTags:   model.TagsFromFilename(filename),
		},
		Shell:    shell,
		Script:   file,
		SubValue: sub,
	}, nil
}

func (s *shell) ShellFor(fs fs.FS, file string) (string, error) {

	var shell, shebang string

	// low-hanging fruit - indicative filename
	if byFile := s.ShellByFilename(file); byFile != "" {
		return byFile, nil
	}

	fd, err := fs.Open(file)
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(fd)
	scanner.Split(bufio.ScanLines)
	if !scanner.Scan() {
		return "", NoContentError
	}
	txt := scanner.Text()
	if strings.HasPrefix(txt, "#!") {
		// shebang found
		for _, potentialShell := range Shells {
			if strings.Contains(txt, potentialShell) {
				shebang = shell
				if loc, err := s.Find(potentialShell); err == nil {
					shell = loc
				} else {
					_, _ = spool.Warn("script has %s but could not find %s (%s)\n", shebang, potentialShell)
				}
			}
		}
	}

	if shebang == "" {
		// if no shebang, just send the first one we find
		for _, s := range Shells {
			if p, err := exec.LookPath(s); err != nil {
				shell = p
			}
		}
	}

	if shell == "" {
		return "", NoShellError
	}

	return shell, nil

}

func (s *shell) ShellByFilename(file string) string {
	ext := filepath.Ext(file)
	if ExtShellMap[ext] != "" {
		sh, err := s.Find(ExtShellMap[ext])
		if err == nil {
			return sh
		}
	}

	return ""
}
