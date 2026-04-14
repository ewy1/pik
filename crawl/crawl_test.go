package crawl

import (
	"github.com/stretchr/testify/assert"
	"io/fs"
	"testing"
)

func TestParentDir(t *testing.T) {
	input := "/var/lib"
	parent := ParentDir(input)
	assert.Equal(t, parent, "/var/")
}

func TestParentDir_TrailingSlash(t *testing.T) {
	input := "/var/lib/"
	parent := ParentDir(input)
	assert.Equal(t, parent, "/var/")
}

func TestParentDir_ToRoot(t *testing.T) {
	input := "/var/"
	parent := ParentDir(input)
	assert.Equal(t, parent, "/")
}
func TestParentDir_ToRoot_NoTrailingSlash(t *testing.T) {
	input := "/var"
	parent := ParentDir(input)
	assert.Equal(t, parent, "/")
}

func TestParentDir_WithoutParent(t *testing.T) {
	input := "/"
	parent := ParentDir(input)
	assert.Equal(t, parent, "/")
}

func TestLocations(t *testing.T) {
	input := "/var/lib/uwu/asdf"
	locs := Locations(input)
	assert.Equal(t, locs, []string{"/var/lib/uwu/asdf", "/var/lib/uwu/", "/var/lib/", "/var/", "/"})
}

func TestLocations_WithDotPath(t *testing.T) {
	input := "/root/./second/asdf/../third"
	locs := Locations(input)
	assert.Equal(t, locs, []string{"/root/second/third", "/root/second/", "/root/", "/"})
}

func TestLocations_HighestFirst(t *testing.T) {
	input := "/one/two/three"
	locs := Locations(input)
	assert.Equal(t, locs[0], "/one/two/three")
}

func Test(t *testing.T) {
	assert.True(t, fs.ValidPath("asdf/hjkl"))
}
