package exc

import (
	"github.com/ewy1/pik/identity"
	"github.com/ewy1/pik/model"
	"github.com/ewy1/pik/runner"
	"github.com/ewy1/pik/spool"
	"io/fs"
	"path/filepath"
)

func (e *exc) Wants(fs fs.FS, file string, entry fs.DirEntry) (bool, error) {
	if entry.IsDir() {
		return false, nil
	}
	info, err := entry.Info()
	if err != nil {
		spool.Warn("%v\n", err)
	}
	return info.Mode()&0100 != 0, nil
}

func (e *exc) CreateTarget(fs fs.FS, source string, file string, entry fs.DirEntry) (model.Target, error) {
	_, filename := filepath.Split(file)
	return &Executable{

		BaseTarget: runner.BaseTarget{
			Identity: identity.New(entry.Name()),
			MyTags:   model.TagsFromFilename(filename),
			MySub:    runner.SubFromFile(file),
		},
		Path: filepath.Join(source, file),
	}, nil
}
