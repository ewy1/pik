package runner

import (
	"pik/identity"
	"pik/model"
)

// BaseTarget is an embeddable type which contains some of the information we need for (almost) every target.
type BaseTarget struct {
	identity.Identity
	MyTags model.Tags
	Sub    []string
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

func (b *BaseTarget) Visible() bool {
	return b.Tags().Visible()
}

func (b *BaseTarget) Invocation(src *model.Source) []string {
	return append([]string{src.Identity.Reduced}, append(b.Sub, b.Identity.Reduced)...)
}

func Hydrated[T model.Target](in T) BaseHydration[T] {
	return BaseHydration[T]{
		Self: in,
	}
}

type BaseHydration[T model.Target] struct {
	Self T
}

func (b *BaseHydration[T]) Icon() string {
	return ""
}

func (b *BaseHydration[T]) Description() string {
	return ""
}

func (b *BaseHydration[T]) Target() model.Target {
	return b.Self
}
