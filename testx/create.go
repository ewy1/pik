//go:build test

package testx

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"pik/identity"
	"pik/model"
	"pik/runner"
	"testing"
)

func TTarget(name string) model.Target {
	return TestTarget{Identifier: name}
}

func TSource(name string, targets ...string) *model.Source {
	src := &model.Source{
		Identity: identity.New(name),
	}
	for _, t := range targets {
		src.Targets = append(src.Targets, TTarget(t))
	}
	return src
}

func TState(sources ...*model.Source) *model.State {
	return &model.State{
		Sources: sources,
	}
}

type TestTarget struct {
	runner.BaseTarget
	Identifier string
	SubValue   []string
	Tags       model.Tags
}

func (t TestTarget) Visible() bool {
	return true
}

func (t TestTarget) Hydrate(src *model.Source) (model.HydratedTarget, error) {
	//TODO implement me
	panic("implement me")
}

func (t TestTarget) Sub() []string {
	return t.SubValue
}

func (t TestTarget) Label() string {
	return t.Identifier
}

func (t TestTarget) Matches(input string) bool {
	return input == t.Identifier
}

func (t TestTarget) Create(s *model.Source) *exec.Cmd {
	panic("whadafak")
}

func AssertTargetIs(t *testing.T, input string, target model.Target) {
	assert.Equal(t, input, target.Label())
}
func AssertTargetIsNot(t *testing.T, input string, target model.Target) {
	assert.NotEqual(t, input, target.Label())
}
func AssertSourceIs(t *testing.T, input string, src *model.Source) {
	assert.Equal(t, input, src.Identity.Reduced)
}
func AssertSourceIsNot(t *testing.T, input string, src *model.Source) {
	assert.NotEqual(t, input, src.Identity.Reduced)
}
