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

func TTarget(name string, sub ...string) model.Target {
	t := TestTarget{Id: identity.New(name), MyTags: model.TagsFromFilename(name), SubValue: sub}
	return &t
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
	runner.Stub
	Id       identity.Identity
	SubValue []string
	MyTags   model.Tags
}

func (t TestTarget) Invocation(src *model.Source) []string {
	return []string{src.Identity.Reduced, t.Id.Reduced}
}

func (t TestTarget) Matches(input string) bool {
	return t.Id.Is(input)
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
	return t.Id.Full
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
	assert.NotNil(t, src.Identity)
	assert.Equal(t, input, src.Identity.Reduced)
}
func AssertSourceIsNot(t *testing.T, input string, src *model.Source) {
	assert.NotNil(t, src.Identity)
	assert.NotEqual(t, input, src.Identity.Reduced)
}
