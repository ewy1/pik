//go:build test

package crawl

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strings"
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

func TestRichLocations(t *testing.T) {
	dir := t.TempDir()
	targetPath := filepath.Join(dir, "target")
	linkPath := filepath.Join(dir, "link")

	innerPath := filepath.Join(targetPath, "inner")
	innerLinkPath := filepath.Join(linkPath, "inner")
	err := os.Mkdir(targetPath, 0777)
	assert.NoError(t, err)
	err = os.Symlink(targetPath, linkPath)
	assert.NoError(t, err)
	err = os.Mkdir(innerPath, 0777)
	result := RichLocations(innerLinkPath)
	assert.NotNil(t, result)
	var hasLink, hasTarget = false, false
	for _, e := range result {
		if strings.HasSuffix(e, "target/inner") {
			hasTarget = true
			continue
		}
		if strings.HasSuffix(e, "link/inner") {
			hasLink = true
			continue
		}
	}
	if !hasLink {
		t.Fatal("evaluated link subdirectories were not included")
	}
	if !hasTarget {
		t.Fatal("original target subdirectories were not included")
	}
}
