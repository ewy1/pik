package js

import (
	"errors"
	"github.com/ewy1/pik/identity"
	"github.com/ewy1/pik/model"
	"github.com/ewy1/pik/runner"
	"os/exec"
	"path/filepath"
	"slices"
)

var jsExtensions = []string{
	".js",
	".cjs",
}

var tsExtensions = []string{
	".ts",
}

var extensions = append(jsExtensions, tsExtensions...)

var managers = []string{
	"pnpm",
	"yarn",
	"npm",
}

var jsInterpreters = []string{
	"node",
	"bun",
}

var tsInterpeters = []string{
	"ts",
	"ts-node",
	"bun",
}

type js struct {
	JsInterpreter string
	TsInterpreter string
	Npm           string
}

var Js = &js{}

var UnsupportedFile = errors.New("unsupported file")
var NoJsInterpreter = errors.New("no js interpreter found in $PATH")
var NoTsInterpreter = errors.New("no ts interpreter found in $PATH")
var NoNpm = errors.New("npm not found in $PATH")

func (n *js) Interpreter(file string) (string, error) {
	ext := filepath.Ext(file)
	if slices.Contains(jsInterpreters, ext) {
		if n.JsInterpreter == "" {
			return "", NoJsInterpreter
		}
		return n.JsInterpreter, nil
	}
	if slices.Contains(tsInterpeters, ext) {
		if n.TsInterpreter == "" {
			return "", NoTsInterpreter
		}
		return n.TsInterpreter, nil
	}
	return "", UnsupportedFile
}

func (n *js) Init() error {
	for _, p := range jsInterpreters {
		if r, err := exec.LookPath(p); err != nil {
			n.JsInterpreter = r
		}
	}
	for _, p := range tsInterpeters {
		if r, err := exec.LookPath(p); err != nil {
			n.TsInterpreter = r
		}
	}

	for _, m := range managers {
		if r, err := exec.LookPath(m); err == nil {
			n.Npm = r
		}
	}

	return nil
}

var npmSub = []string{
	"npm",
}

func (n *js) CreateRun(name, cmd string) model.Target {
	return &Npm{
		BaseTarget: runner.BaseTarget{
			Identity: identity.New(cmd),
			MyTags:   model.TagsFromFilename(cmd),
			MySub:    npmSub,
		},
		Name: name,
		Cmd:  cmd,
	}
}
