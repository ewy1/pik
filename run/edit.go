package run

import (
	"errors"
	"os"
	"os/exec"
	"pik/model"
)

var NoEditorError = errors.New("$EDITOR not set")

func Edit(t model.Target, src *model.Source) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return NoEditorError
	}
	cmd := exec.Command(editor, t.File(src))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = src.Path
	return cmd.Run()
}
