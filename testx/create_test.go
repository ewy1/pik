//go:build test

package testx

import (
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
