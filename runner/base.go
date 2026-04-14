package runner

import (
	"os/exec"
	"pik/identity"
	"pik/model"
)

type BaseTarget struct {
	identity.Identity
	MyTags model.Tags
}

func (t *BaseTarget) Tags() model.Tags {
	return t.MyTags
}

func (t *BaseTarget) Matches(input string) bool {
	return t.Identity.Is(input)
}

func (t *BaseTarget) ShortestId() string {
	return t.Reduced
}

func Hydrated[T model.Target](in T) BaseHydration[T] {
	return BaseHydration[T]{
		BaseTarget: in,
	}
}

type BaseHydration[T model.Target] struct {
	BaseTarget T
}

func (b BaseHydration[T]) Matches(input string) bool {
	return b.BaseTarget.Matches(input)
}

func (b BaseHydration[T]) Create(s *model.Source) *exec.Cmd {
	return b.BaseTarget.Create(s)
}

func (b BaseHydration[T]) Sub() []string {
	return b.BaseTarget.Sub()
}

func (b BaseHydration[T]) Label() string {
	return b.BaseTarget.Label()
}

func (b BaseHydration[T]) Hydrate(src *model.Source) (model.HydratedTarget, error) {
	return b, nil
}

func (b BaseHydration[T]) Visible() bool {
	return b.BaseTarget.Visible()
}

func (b BaseHydration[T]) Tags() model.Tags {
	return b.BaseTarget.Tags()
}

func (b BaseHydration[T]) ShortestId() string {
	return b.BaseTarget.ShortestId()
}

func (b BaseHydration[T]) Icon() string {
	return "  "
}

func (b BaseHydration[T]) Description() string {
	return ""
}

func (b BaseHydration[T]) Target() model.Target {
	return b.BaseTarget
}

func (b BaseTarget) Visible() bool {
	return b.Tags().Visible()
}
