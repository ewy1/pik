//go:build test

package identity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentity(t *testing.T) {
	a := New("asdf")
	b := New("asdf")
	c := New("hhhh")
	assert.True(t, a.I(b))
	assert.False(t, a.I(c))
}

func TestIdentity_Extension(t *testing.T) {
	a := New("asdf.sh")
	b := New("asdf")
	c := New("hhhh")
	assert.True(t, a.I(b))
	assert.False(t, a.I(c))
}

func TestIdentity_Hidden(t *testing.T) {
	a := New(".asdf")
	b := New("asdf")
	c := New("hhhh")
	assert.True(t, a.I(b))
	assert.False(t, a.I(c))
}

func TestIdentity_I_Tagged(t *testing.T) {
	a := New("asdf.pre.sh")
	b := New("asdf")
	c := New("hhhh")
	assert.True(t, a.I(b))
	assert.False(t, a.I(c))
}

func TestIdentity_Is(t *testing.T) {
	a := New(".asdf.iahsodiu.txt")
	b := "asdf"
	assert.True(t, a.Is(b))
}
