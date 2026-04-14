package search

import (
	"github.com/stretchr/testify/assert"
	"pik/testx"
	"testing"
)

func TestSearch_TargetOnly(t *testing.T) {
	st := testx.TState(testx.TSource("src", "abc", "def"))
	target, source, _, _, _ := Search(st, "def")
	testx.AssertSourceIs(t, "src", source)
	testx.AssertTargetIs(t, "def", target)
}

func TestSearch_TargetAndSource(t *testing.T) {
	st := testx.TState(testx.TSource("src", "abc", "def"))
	target, source, _, _, _ := Search(st, "src", "def")
	testx.AssertSourceIs(t, "src", source)
	testx.AssertTargetIs(t, "def", target)
}

func TestSearch_SourceDefaultTarget(t *testing.T) {
	st := testx.TState(testx.TSource("src", "abc", "src"))
	target, src, _, _, _ := Search(st, "src")
	testx.AssertSourceIs(t, "src", src)
	assert.NotNil(t, target)
}

func TestSearch_SubdirWrong(t *testing.T) {
	st := testx.TState(testx.TSource("src", "abc", "src"))
	st.Sources[0].Targets = append(st.Sources[0].Targets, testx.TestTarget{
		Identifier: "script",
		SubValue:   []string{"subdir"},
	})
	target, src, confirm, sd, _ := Search(st, "wrong", "script")
	testx.AssertSourceIs(t, "src", src)
	testx.AssertTargetIs(t, "script", target)
	assert.Equal(t, sd, []string{"wrong"})
	assert.NotNil(t, target)
	assert.True(t, confirm)
}

func TestSearch_SubdirMissing(t *testing.T) {
	st := testx.TState(testx.TSource("src", "abc", "src"))
	st.Sources[0].Targets = append(st.Sources[0].Targets, testx.TestTarget{
		Identifier: "script",
		SubValue:   []string{"subdir"},
	})
	target, src, confirm, sd, _ := Search(st, "script")
	testx.AssertSourceIs(t, "src", src)
	testx.AssertTargetIs(t, "script", target)
	assert.Nil(t, sd)
	assert.NotNil(t, target)
	assert.False(t, confirm)
}

func TestSearch_Args(t *testing.T) {
	st := testx.TState(testx.TSource("src", "abc", "def"))
	target, source, _, _, args := Search(st, "def", "a1", "a2")
	testx.AssertSourceIs(t, "src", source)
	testx.AssertTargetIs(t, "def", target)
	assert.Equal(t, []string{"a1", "a2"}, args)
}

func TestSearch_Args_SubdirMissing(t *testing.T) {
	st := testx.TState(testx.TSource("src", "abc", "src"))
	st.Sources[0].Targets = append(st.Sources[0].Targets, testx.TestTarget{
		Identifier: "script",
		SubValue:   []string{"subdir"},
	})
	target, src, _, _, args := Search(st, "script", "a1", "a2")
	testx.AssertSourceIs(t, "src", src)
	testx.AssertTargetIs(t, "script", target)
	assert.Equal(t, []string{"a1", "a2"}, args)
}

func TestSearch_Args_SubdirPresent(t *testing.T) {
	st := testx.TState(testx.TSource("src", "abc", "src"))
	st.Sources[0].Targets = append(st.Sources[0].Targets, testx.TestTarget{
		Identifier: "script",
		SubValue:   []string{"subdir"},
	})
	target, src, _, _, args := Search(st, "subdir", "script", "a1", "a2")
	testx.AssertSourceIs(t, "src", src)
	testx.AssertTargetIs(t, "script", target)
	assert.Equal(t, []string{"a1", "a2"}, args)
}

func TestSearch_SecondarySource(t *testing.T) {
	st := testx.TState(testx.TSource("src", "abc", "def"), testx.TSource("aaa", "hjkl"))
	target, source, _, _, _ := Search(st, "aaa", "hjkl")
	testx.AssertSourceIs(t, "aaa", source)
	testx.AssertTargetIs(t, "hjkl", target)
}

func TestSearch_SecondarySource_DuplicateTargetName(t *testing.T) {
	st := testx.TState(testx.TSource("src", "abc", "def"), testx.TSource("aaa", "abc"))
	target, source, confirm, _, _ := Search(st, "aaa", "def")
	testx.AssertSourceIs(t, "src", source)
	testx.AssertTargetIs(t, "def", target)
	assert.True(t, confirm)
}

func TestSearch_SourceTargetMixup(t *testing.T) {
	st := testx.TState(testx.TSource("src", "abc"), testx.TSource("aaa", "ccc"))
	target, source, confirm, _, _ := Search(st, "src", "ccc")
	testx.AssertSourceIs(t, "aaa", source)
	testx.AssertTargetIs(t, "ccc", target)
	assert.True(t, confirm)
}
