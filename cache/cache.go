package cache

import (
	"bufio"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"pik/model"
	"pik/paths"
	"strings"
)

type Cache struct {
	Entries []Entry
}

// Merge combines two caches and filters duplicate keys
func (c Cache) Merge(other Cache) Cache {
	mp := make(map[string]string)
	for _, e := range append(c.Entries, other.Entries...) {
		mp[e.Path] = e.Label
	}
	result := Cache{}
	for p, l := range mp {
		result.Entries = append(result.Entries, Entry{Label: l, Path: p})
	}
	return result
}

type Entry struct {
	Path  string
	Label string
}

var Empty = Cache{}

// Path is the file path to the "contexts" cache file
var Path = path.Join(paths.Cache, "contexts")

// FsPath is the Path with the leading slash removed, to be opened from fs.FS
var FsPath = Path[1:]

var UnexpectedEntryError = errors.New("unexpected cache entry")

// LoadFile creates a Cache from a file or an empty one if the file does not exist
// this handles opening a reader for Unmarshal
func LoadFile(root fs.FS, path string) (Cache, error) {
	fd, err := root.Open(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Cache{}, nil
	} else if err != nil {
		return Cache{}, err
	}
	if fd != nil {
		defer fd.Close()
	}
	return Unmarshal(fd)
}

// Unmarshal attempts to create a Cache from reader content
func Unmarshal(r io.Reader) (Cache, error) {
	c := Cache{}
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
		default:
			return c, UnexpectedEntryError
		}
		c.Entries = append(c.Entries, *entry)
	}
	return c, nil
}

// Marshal returns the file representation of the Cache
func (c Cache) Marshal() []byte {
	b := strings.Builder{}
	for _, e := range c.Entries {
		b.WriteString(e.Path)
		b.WriteString(" # ")
		b.WriteString(e.Label)
		b.WriteString("\n")
	}
	return []byte(b.String())
}

func (c Cache) String() string {
	return string(c.Marshal())
}

func New(st *model.State) Cache {
	c := &Cache{}
	for _, s := range st.Sources {
		c.Entries = append(c.Entries, Entry{
			Path:  s.Path,
			Label: s.Label(),
		})
	}
	return *c
}

func SaveFile(path string, s *model.State, loaded Cache) error {
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	if fd != nil {
		defer fd.Close()
	}
	return Save(s, fd, loaded)
}

func Save(s *model.State, w io.Writer, loaded Cache) error {
	result := New(s).Merge(loaded)
	_, err := w.Write([]byte(result.Marshal()))
	return err

}

func LoadState(f fs.FS, cache Cache, indexers []model.Indexer, runners []model.Runner) (*model.State, []error) {
	var locs []string
	for _, e := range cache.Entries {
		locs = append(locs, e.Path)
	}
	return model.NewState(f, locs, indexers, runners)
}

func (c Cache) Strip(needle Cache) Cache {
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
