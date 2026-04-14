//go:build test

package model

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

type taggedTarget struct {
	MyTags Tags
}

func (t taggedTarget) Matches(input string) bool {
	//TODO implement me
	panic("implement me")
}

func (t taggedTarget) Create(s *Source) *exec.Cmd {
	//TODO implement me
	panic("implement me")
}

func (t taggedTarget) Sub() []string {
	//TODO implement me
	panic("implement me")
}

func (t taggedTarget) Label() string {
	//TODO implement me
	panic("implement me")
}

func (t taggedTarget) Hydrate(src *Source) (HydratedTarget, error) {
	//TODO implement me
	panic("implement me")
}

func (t taggedTarget) Tags() Tags {
	return t.MyTags
}

func (t taggedTarget) ShortestId() string {
	//TODO implement me
	panic("implement me")
}

func tagged(in ...Tag) taggedTarget {
	return taggedTarget{MyTags: tags(in...)}
}

func tags(in ...Tag) Tags {
	return Tags(in)
}

func TestTags_Count(t *testing.T) {
	input := tags(Here, Hidden)
	assert.Len(t, input, 2)
}

func TestTags_Is(t *testing.T) {
	input := tags(Final, Post)
	assert.True(t, input.Has(Final))
	assert.True(t, input.Has(Post))
	assert.False(t, input.Has(Pre))
}

func TestTags_AnyOf(t *testing.T) {
	input := tags(Final, Post)
	assert.True(t, input.AnyOf(Final))
	assert.False(t, input.AnyOf(Pre, Hidden))
}

func TestTags_AnyOf_Mix(t *testing.T) {
	input := tags(Final, Post)
	assert.True(t, input.AnyOf(Post, Hidden))
}

func TestTags_AnyOf_EmptyInput(t *testing.T) {
	input := tags()
	assert.False(t, input.AnyOf(Final))
	assert.False(t, input.AnyOf(Pre, Hidden))
}

func TestTags_AnyOf_EmptySearch(t *testing.T) {
	input := tags(Final, Post)
	assert.True(t, input.AnyOf())
}

func TestTagsFromFilename(t *testing.T) {
	inputs := []any{
		*Pre,
		*Here,
		*Final,
	}
	input := fmt.Sprintf("script.%v.%v.%v.ext", inputs...)
	output := TagsFromFilename(input)
	assert.Len(t, output, len(inputs))
	assert.Contains(t, output, Pre)
	assert.Contains(t, output, Here)
	assert.Contains(t, output, Final)
}

func TestTagsHidden_VisibleByDefault(t *testing.T) {
	input := tags(Here)
	target := tagged(input...)
	assert.True(t, target.Tags().Visible())
}

func TestTagsHidden_FromHidden(t *testing.T) {
	input := tags(Hidden)
	target := tagged(input...)
	assert.False(t, target.Tags().Visible())
}

func TestTagsHidden_FromFilename(t *testing.T) {
	input := TagsFromFilename(".asdf.sh")
	target := tagged(input...)
	assert.False(t, target.Tags().Visible())
}

func TestTagsHidden_FromTrigger(t *testing.T) {
	input := tags(Pre, Here)
	target := tagged(input...)
	assert.False(t, target.Tags().Visible())
}
