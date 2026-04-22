//go:build test

package testx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAssertSourceIs_Correct(t *testing.T) {
	src := TSource("abc", "def")
	AssertSourceIs(t, "abc", src)
}

func TestAssertSourceIs_Wrong(t *testing.T) {
	src := TSource("abc", "def")
	AssertSourceIsNot(t, ";lkjh", src)
}

func TestAssertTargetIs_Correct(t *testing.T) {
	ta := TTarget("aaaa")
	AssertTargetIs(t, "aaaa", ta)
}

func TestAssertTargetIs_Wrong(t *testing.T) {
	ta := TTarget("aaaa")
	AssertTargetIsNot(t, "bbbbbb", ta)
}

func TestTTargetIdentity(t *testing.T) {
	ta := TTarget("asdf.hidden.sh")
	assert.True(t, ta.Matches("asdf"))
	assert.False(t, ta.Matches("hidden"))
}
