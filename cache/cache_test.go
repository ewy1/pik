//go:build test

package cache

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	. "pik/testx"
	"strings"
	"testing"
	"testing/fstest"
)

func TestLoadFile(t *testing.T) {
	input := `

# comment
/asdf/hjkl # label

# other comment
//comment

`
	f := fstest.MapFS{
		"contexts": &fstest.MapFile{
			Data: []byte(input),
		},
	}
	res, err := LoadFile(f, "contexts")
	assert.NoError(t, err)
	assert.Len(t, res.Entries, 1)
	assert.Equal(t, res.Entries[0].Path, "/asdf/hjkl")
	assert.Equal(t, res.Entries[0].Label, "label")
}

func TestFromReader_Blank(t *testing.T) {
	input := `   
`
	sr := strings.NewReader(input)
	c, err := Unmarshal(sr)
	assert.Nil(t, err)
	assert.Len(t, c.Entries, 0)
}

func TestFromReader_OneEntry(t *testing.T) {
	input := `/abc/def # deffers`
	sr := strings.NewReader(input)
	c, err := Unmarshal(sr)
	assert.Nil(t, err)
	assert.Len(t, c.Entries, 1)
	assert.Equal(t, c.Entries[0], Entry{
		Path:  "/abc/def",
		Label: "deffers",
	})
}

func TestFromReader_ManyEntries(t *testing.T) {
	input := `/abc/def # deffers
/123/aa # i love aa
/path/src # da source
`
	sr := strings.NewReader(input)
	c, err := Unmarshal(sr)
	assert.Nil(t, err)
	assert.Len(t, c.Entries, 3)
	assert.Equal(t, c.Entries[0], Entry{
		Path:  "/abc/def",
		Label: "deffers",
	})
	assert.Equal(t, c.Entries[1], Entry{
		Path:  "/123/aa",
		Label: "i love aa",
	})
	assert.Equal(t, c.Entries[2], Entry{
		Path:  "/path/src",
		Label: "da source",
	})
}

func TestFromReader_Comments(t *testing.T) {
	input := `
// comment
/abc/def # deffers
# comment
/123/aa # i love aa
// # comment
/path/src # da source
# // comment
`
	sr := strings.NewReader(input)
	c, err := Unmarshal(sr)
	assert.Nil(t, err)
	assert.Len(t, c.Entries, 3)
	assert.Equal(t, c.Entries[0], Entry{
		Path:  "/abc/def",
		Label: "deffers",
	})
	assert.Equal(t, c.Entries[1], Entry{
		Path:  "/123/aa",
		Label: "i love aa",
	})
	assert.Equal(t, c.Entries[2], Entry{
		Path:  "/path/src",
		Label: "da source",
	})
}

func TestStrip(t *testing.T) {
	c := Cache{Entries: []Entry{{"/asdf/123", ""}, {"xxxxx", "lab"}}}
	remove := Cache{Entries: []Entry{{"xxxxx", "wronglabel"}}}
	result := c.Strip(remove)
	assert.Equal(t, Cache{Entries: []Entry{{"/asdf/123", ""}}}, result)
}

func TestStrip_Nothing(t *testing.T) {
	c := Cache{Entries: []Entry{{"/asdf/123", ""}, {"/asdf/123", ""}}}
	old := Cache{}
	result := c.Strip(old)
	assert.Equal(t, c, result)
}

func TestMerge(t *testing.T) {
	a := Entry{
		Path: "/usr/share/asdf",
	}
	b := Entry{
		Path: "/test/location",
	}
	c := Entry{
		Path:  "/new/mypath",
		Label: "mypath",
	}
	base := Cache{Entries: []Entry{
		a, b,
	}}
	other := Cache{Entries: []Entry{
		b, c,
	}}
	result := base.Merge(other)
	assert.Len(t, result.Entries, 3)
	assert.Contains(t, result.Entries, a, b, c)
}

func TestCache_String(t *testing.T) {
	c := Cache{Entries: []Entry{{Path: "/asdf/hjkl", Label: "hjkl"}}}
	output := c.String()
	for _, e := range c.Entries {
		assert.Contains(t, output, e.Path)
		assert.Contains(t, output, e.Label)
	}
}

func TestNew(t *testing.T) {
	st := TState(TSource("/test/location", "target1", "target2"), TSource("/other/place/asdf/hjkl", "1target", "2target"))
	c := New(st)
	assert.NotNil(t, c)
	assert.Len(t, c.Entries, 2)
	assert.Equal(t, "/test/location", c.Entries[0].Path)
	assert.Equal(t, "/other/place/asdf/hjkl", c.Entries[1].Path)
}

func TestLoadFile_NotExist(t *testing.T) {
	f := fstest.MapFS{}
	c, err := LoadFile(f, "anything is fine")
	assert.Nil(t, c.Entries)
	assert.NoError(t, err)

}

func TestSaveFile(t *testing.T) {
	dir := t.TempDir()
	loc := filepath.Join(dir, "savefile")
	st := TState(TSource("source_one", "target"), TSource("second_source", "t1", "t2", "t3"))
	c := Cache{Entries: []Entry{
		{
			Path:  "path",
			Label: "label",
		},
		{
			Path:  "/otherpath/123",
			Label: "",
		},
	}}
	err := SaveFile(loc, st, c)
	assert.NoError(t, err)
}
