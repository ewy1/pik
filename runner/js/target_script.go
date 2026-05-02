package js

import (
	"github.com/ewy1/pik/describe"
	"github.com/ewy1/pik/model"
	"github.com/ewy1/pik/runner"
	"github.com/ewy1/pik/spool"
	"os/exec"
	"path/filepath"
)

type Script struct {
	runner.BaseTarget
	Typed bool
}

func (t *Script) Icon() string {
	if t.Typed {
		return "\uE628"
	} else {
		return "\uE60C"
	}
}

func (t *Script) Description(src *model.HydratedSource) string {
	d, err := describe.Describe(t, t.File(src.Source))
	if err != nil {
		_, _ = spool.Warn("%v\n", err)
	}
	return d
}

func (t *Script) Target() model.Target {
	return t
}

func (t *Script) Create(s *model.Source) *exec.Cmd {
	if t.Typed {
		return exec.Command(Js.TsInterpreter, t.File(s))
	} else {
		return exec.Command(Js.JsInterpreter, t.File(s))
	}
}

func (t *Script) Label() string {
	return t.Identity.Full
}

func (t *Script) Hydrate(src *model.Source) (model.HydratedTarget, error) {
	return t, nil
}

func (t *Script) File(src *model.Source) string {
	return filepath.Join(src.Path, "package.json")
}
