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

var Path = path.Join(paths.Cache, "contexts")

var UnexpectedEntryError = errors.New("unexpected cache entry")

func Load() (Cache, error) {
	fd, err := os.Open(Path)
	if errors.Is(err, os.ErrNotExist) {
		return Cache{}, nil
	} else if err != nil {
		return Cache{}, err
	}
	defer fd.Close()
	return FromReader(fd)
}

func FromReader(r io.Reader) (Cache, error) {
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

func (c Cache) String() string {
	b := strings.Builder{}
	for _, e := range c.Entries {
		b.WriteString(e.Path)
		b.WriteString(" # ")
		b.WriteString(e.Label)
		b.WriteString("\n")
	}
	return b.String()
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

func Save(s *model.State) error {
	ld, err := Load()
	if err != nil {
		return err
	}
	c := New(s).Merge(ld)
	return os.WriteFile(Path, []byte(c.String()), os.ModePerm)
}

func LoadState(f fs.FS, cache Cache, indexers []model.Indexer, runners []model.Runner) (*model.State, error) {
	var locs []string
	for _, e := range cache.Entries {
		locs = append(locs, e.Path)
	}
	return model.NewState(f, locs, indexers, runners)
}
