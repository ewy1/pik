package pikdex

import (
	"github.com/ewy1/pik/paths"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"testing"
	"testing/fstest"
)

func TestUwudex_WantsWalk_AnyRoot(t *testing.T) {
	for _, r := range paths.Roots {
		data := fstest.MapFS{
			r: &fstest.MapFile{
				Data: nil,
				Mode: fs.ModeDir,
			},
		}
		u := &pikdex{}
		result, root, err := u.WantsWalk(data)
		assert.Equal(t, root, r)
		assert.NoError(t, err)
		assert.True(t, result)
	}
}

func TestUwudex_WantsWalk_TwoRoots(t *testing.T) {
	data := fstest.MapFS{
		paths.Roots[0]: &fstest.MapFile{
			Data: nil,
			Mode: fs.ModeDir,
		},
		paths.Roots[1]: &fstest.MapFile{
			Data: nil,
			Mode: fs.ModeDir,
		},
	}
	u := &pikdex{}
	result, r, err := u.WantsWalk(data)
	// no guarantee we pick any one lol
	assert.Contains(t, paths.Roots, r)
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestUwudex_WantsWalk_NoRoots(t *testing.T) {
	data := fstest.MapFS{
		"asdf.txt": &fstest.MapFile{},
	}
	u := &pikdex{}
	result, r, err := u.WantsWalk(data)
	assert.Equal(t, "", r)
	assert.NoError(t, err)
	assert.False(t, result)
}
