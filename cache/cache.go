package cache

import (
	"bufio"
	"errors"
	"github.com/ewy1/pik/model"
	"github.com/ewy1/pik/paths"
	"io"
	"io/fs"
	"os"
	"slices"
	"strings"
)

type Cache struct {
	Entries []Entry
}

type cacheInit struct{}

var Init model.Initializer = &cacheInit{}

func (i *cacheInit) Init() error {
	return nil
}

// Merge combines two caches and filters duplicate keys
func (c *Cache) Merge(other *Cache) *Cache {
	switch {
	case other == nil && c != nil:
		return c
	case c == nil && other != nil:
		return other
	case c == nil:
		return nil

	}
	mp := make(map[string]string)
	for _, e := range append(c.Entries, other.Entries...) {
		mp[e.Path] = e.Label
	}
	result := &Cache{}
	for p, l := range mp {
		result.Entries = append(result.Entries, Entry{Label: l, Path: p})
	}
	return result
}

type Entry struct {
	Path  string
	Label string
}

func (e Entry) String() string {
	return e.Path + " # " + e.Label
}

// LoadFile creates a Cache from a file or an empty one if the file does not exist
// this handles opening a reader for Unmarshal
func LoadFile(root fs.FS, path string) (*Cache, error) {
	fd, err := root.Open(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if fd != nil {
		defer fd.Close()
	}
	return Unmarshal(fd)
}

// Unmarshal attempts to create a Cache from reader content
func Unmarshal(r io.Reader) (*Cache, error) {
	c := &Cache{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == '#' || line[0:2] == "//" {
			continue
		}

		entry := &Entry{}
		parts := strings.SplitN(line, "#", 2)
		switch len(parts) {
		case 2:
			entry.Label = strings.TrimSpace(parts[1])
			fallthrough
		case 1:
			entry.Path = strings.TrimSpace(parts[0])
		}
		c.Entries = append(c.Entries, *entry)
	}
	return c, nil
}

// Marshal returns the file representation of the Cache
func (c *Cache) Marshal() []byte {
	b := strings.Builder{}
	for _, e := range c.Entries {
		b.WriteString(e.String())
		b.WriteString("\n")
	}
	return []byte(b.String())
}

func (c *Cache) String() string {
	return string(c.Marshal())
}
func New(st *model.State) *Cache {
	c := &Cache{}
	for _, s := range st.Sources {
		c.Entries = append(c.Entries, Entry{
			Path:  s.Path,
			Label: s.Label(),
		})
	}
	return c
}

func MergeAndSave(in *model.State) error {
	root := "/"
	f := os.DirFS(root)
	// remove leading slash from the dirfs
	loaded, err := LoadFile(f, strings.TrimPrefix(paths.ContextsFile.String(), "/"))
	if err != nil {
		return err
	}
	insert := New(in)
	result := loaded.Merge(insert)
	if loaded == nil {
		return SaveFile(paths.ContextsFile.String(), result)
	}
	if slices.Equal(loaded.Entries, result.Entries) {
		return nil
	}
	return SaveFile(paths.ContextsFile.String(), result)
}

// SaveFile helps you use Save with a file path instead of a reader
func SaveFile(path string, loaded *Cache) error {
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	if fd != nil {
		defer fd.Close()
	}
	return Save(fd, loaded)
}

// Save writes a cache to the contexts file
func Save(w io.Writer, loaded *Cache) error {
	_, err := w.Write(loaded.Marshal())
	return err

}

// LoadState creates a state with model.NewState based on cache content
func LoadState(f fs.FS, cache *Cache, indexers []model.Indexer, runners []model.Runner) (*model.State, []error) {
	var locs []string
	for _, e := range cache.Entries {
		locs = append(locs, e.Path)
	}
	return model.NewState(f, locs, indexers, runners)
}

// Strip removes the needle's entries from the receiver's entries when they have matching paths.
// used to skip already indexed locations when auto-all-ing
func (c *Cache) Strip(needle Cache) Cache {
	var result []Entry
outer:
	for _, e := range c.Entries {
		for _, t := range needle.Entries {
			if t.Path == e.Path {
				continue outer
			}
		}
		result = append(result, e)
	}
	return Cache{
		Entries: result,
	}
}
