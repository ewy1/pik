//go:build test

package runner

import (
	"os/exec"
	"pik/model"
)

// Stub is the most minimal and useless implementation of the target interface. It only panics. Use if you need a target-compliant struct.
type Stub struct {
}

func (s Stub) Matches(input string) bool {
	//TODO implement me
	panic("implement me")
}

func (s Stub) Create(src *model.Source) *exec.Cmd {
	//TODO implement me
	panic("implement me")
}

func (s Stub) Sub() []string {
	//TODO implement me
	panic("implement me")
}

func (s Stub) Label() string {
	//TODO implement me
	panic("implement me")
}

func (s Stub) Hydrate(src *model.Source) (model.HydratedTarget, error) {
	//TODO implement me
	panic("implement me")
}

func (s Stub) Tags() model.Tags {
	return nil
}

func (s Stub) ShortestId() string {
	//TODO implement me
	panic("implement me")
}

func (s Stub) Visible() bool {
	//TODO implement me
	panic("implement me")
}

func (s Stub) Invocation(src *model.Source) []string {
	//TODO implement me
	panic("implement me")
}
