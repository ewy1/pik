//go:build test

package pikdex

import (
	"github.com/stretchr/testify/assert"
	"pik/testx"
	"testing"
)

func TestHydrate(t *testing.T) {
	aliases := []string{"alias1", "alias2"}
	p := pikdex{
		mods: map[string]*SourceData{
			"asdf": {
				Aliases: aliases,
				Icon:    "I",
				Path:    "asdf",
			},
		},
	}
	src := testx.TSource("asdf", "target")
	hyd := src.Hydrate(nil)
	assert.NotNil(t, hyd)
	err := p.Mod(src, hyd)
	assert.NoError(t, err)
	assert.Equal(t, aliases, hyd.Aliases)
	assert.Equal(t, "I", hyd.Icon)
}
