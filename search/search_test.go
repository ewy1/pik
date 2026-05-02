//go:build test

package search

import (
	"git.ewy.one/pik.git/runner"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearch_TargetOnly(t *testing.T) {
	st := runner.TState(runner.TSource("src", "abc", "def"))
	res := Search(st, "def")
	runner.AssertSourceIs(t, "src", res.Source)
	runner.AssertTargetIs(t, "def", res.Target)
}

func TestSearch_TargetAndSource(t *testing.T) {
	st := runner.TState(runner.TSource("src", "abc", "def"))
	res := Search(st, "src", "def")
	runner.AssertSourceIs(t, "src", res.Source)
	runner.AssertTargetIs(t, "def", res.Target)
}

func TestSearch_TargetAndSource_CaseInsensitive(t *testing.T) {
	st := runner.TState(runner.TSource("src", "abc", "def"))
	res := Search(st, "SRC", "DeF")
	runner.AssertSourceIs(t, "src", res.Source)
	runner.AssertTargetIs(t, "def", res.Target)
}

func TestSearch_SourceDefaultTarget(t *testing.T) {
	st := runner.TState(runner.TSource("src", "abc", "src"))
	res := Search(st, "src")
	runner.AssertSourceIs(t, "src", res.Source)
	assert.NotNil(t, res.Target)
}

func TestSearch_SubdirWrong(t *testing.T) {
	st := runner.TState(runner.TSource("src", "abc", "src"))
	st.Sources[0].Targets = append(st.Sources[0].Targets, runner.TTarget("script", "subdir"))
	res := Search(st, "wrong", "script")
	runner.AssertSourceIs(t, "src", res.Source)
	runner.AssertTargetIs(t, "script", res.Target)
	assert.Equal(t, []string{"wrong"}, res.Sub)
	assert.NotNil(t, res.Target)
	assert.True(t, res.NeedsConfirmation)
}

func TestSearch_SubdirMissing(t *testing.T) {
	st := runner.TState(runner.TSource("src", "abc", "src"))
	st.Sources[0].Targets = append(st.Sources[0].Targets, runner.TTarget("script", "subdir"))
	res := Search(st, "script")
	runner.AssertSourceIs(t, "src", res.Source)
	runner.AssertTargetIs(t, "script", res.Target)
	assert.Nil(t, res.Sub)
	assert.NotNil(t, res.Target)
	assert.False(t, res.NeedsConfirmation)
}

func TestSearch_Args(t *testing.T) {
	st := runner.TState(runner.TSource("src", "abc", "def"))
	res := Search(st, "def", "a1", "a2")
	runner.AssertSourceIs(t, "src", res.Source)
	runner.AssertTargetIs(t, "def", res.Target)
	assert.Equal(t, []string{"a1", "a2"}, res.Args)
}

func TestSearch_Args_SubdirMissing(t *testing.T) {
	st := runner.TState(runner.TSource("src", "abc", "src"))
	st.Sources[0].Targets = append(st.Sources[0].Targets, runner.TTarget("script", "subdir"))
	res := Search(st, "script", "a1", "a2")
	runner.AssertSourceIs(t, "src", res.Source)
	runner.AssertTargetIs(t, "script", res.Target)
	assert.Equal(t, []string{"a1", "a2"}, res.Args)
}

func TestSearch_Args_SubdirPresent(t *testing.T) {
	st := runner.TState(runner.TSource("src", "abc", "src"))
	st.Sources[0].Targets = append(st.Sources[0].Targets, runner.TTarget("script", "subdir"))
	res := Search(st, "subdir", "script", "a1", "a2")
	runner.AssertSourceIs(t, "src", res.Source)
	runner.AssertTargetIs(t, "script", res.Target)
	assert.Equal(t, []string{"a1", "a2"}, res.Args)
}

func TestSearch_SecondarySource(t *testing.T) {
	st := runner.TState(runner.TSource("src", "abc", "def"), runner.TSource("aaa", "hjkl"))
	res := Search(st, "aaa", "hjkl")
	runner.AssertSourceIs(t, "aaa", res.Source)
	runner.AssertTargetIs(t, "hjkl", res.Target)
}

func TestSearch_SecondarySource_DuplicateTargetName(t *testing.T) {
	st := runner.TState(runner.TSource("src", "abc", "def"), runner.TSource("aaa", "abc"))
	res := Search(st, "aaa", "def")
	runner.AssertSourceIs(t, "src", res.Source)
	runner.AssertTargetIs(t, "def", res.Target)
	assert.True(t, res.NeedsConfirmation)
}

func TestSearch_SourceTargetMixup(t *testing.T) {
	st := runner.TState(runner.TSource("src", "abc"), runner.TSource("aaa", "ccc"))
	res := Search(st, "src", "ccc")
	runner.AssertSourceIs(t, "aaa", res.Source)
	runner.AssertTargetIs(t, "ccc", res.Target)
	assert.True(t, res.NeedsConfirmation)
}

func TestSearch_Override(t *testing.T) {
	st := runner.TState(runner.TSource("src", "abc.override.sh", "abc.sh"))
	res := Search(st, "src", "abc")
	assert.Equal(t, "abc.override.sh", res.Target.(*runner.TestTarget).Id.Full)
	assert.False(t, res.NeedsConfirmation)
}

func TestSearch_SubdirDefault(t *testing.T) {
	tgt := runner.TTarget("subname", "subname")
	src := runner.TSource("src")
	src.Targets = append(src.Targets, tgt)
	st := runner.TState(src)
	res := Search(st, "subname")
	assert.Nil(t, res.Args)
	assert.Equal(t, res.Target, tgt)
	assert.False(t, res.NeedsConfirmation)
}

func TestSearch_SourceDefault(t *testing.T) {
	src := runner.TSource("sourcename", "sourcename")
	st := runner.TState(src)
	res := Search(st, "sourcename")
	assert.Nil(t, res.Args)
	assert.Equal(t, res.Target, src.Targets[0])
	assert.False(t, res.NeedsConfirmation)
}

func TestSearch_SourceDefault_Other(t *testing.T) {
	src := runner.TSource("src", "src", "other")
	st := runner.TState(src)
	res := Search(st, "src", "other")
	assert.Nil(t, res.Args)
	assert.Equal(t, res.Target, src.Targets[1])
	assert.False(t, res.NeedsConfirmation)
}

func TestSearch_SubdirDefault_Other(t *testing.T) {
	tgt := runner.TTarget("subname", "subname")
	other := runner.TTarget("othername", "subname")
	src := runner.TSource("src")
	src.Targets = append(src.Targets, tgt, other)
	st := runner.TState(src)
	res := Search(st, "subname", "othername")
	assert.Nil(t, res.Args)
	assert.Equal(t, other, res.Target)
	assert.False(t, res.NeedsConfirmation)
}
