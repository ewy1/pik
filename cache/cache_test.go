package cache

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestFromReader_Blank(t *testing.T) {
	input := `   
`
	sr := strings.NewReader(input)
	c, err := FromReader(sr)
	assert.Nil(t, err)
	assert.Len(t, c.Entries, 0)
}

func TestFromReader_OneEntry(t *testing.T) {
	input := `/abc/def # deffers`
	sr := strings.NewReader(input)
	c, err := FromReader(sr)
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
	c, err := FromReader(sr)
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
	c, err := FromReader(sr)
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
