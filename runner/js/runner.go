package js

import (
	"github.com/ewy1/pik/identity"
	"github.com/ewy1/pik/model"
	"github.com/ewy1/pik/runner"
	"io/fs"
	"path/filepath"
	"slices"
)

func (n *js) Wants(fs fs.FS, file string, entry fs.DirEntry) (bool, error) {
	ext := filepath.Ext(entry.Name())
	return slices.Contains(extensions, ext), nil
}

func (n *js) CreateTarget(fs fs.FS, source string, file string, entry fs.DirEntry) (model.Target, error) {
	ext := filepath.Ext(entry.Name())
	typed := slices.Contains(tsExtensions, ext)
	return &Script{
		BaseTarget: runner.BaseTarget{
			Identity: identity.New(entry.Name()),
			MyTags:   model.TagsFromFilename(entry.Name()),
			MySub:    runner.SubFromFile(file),
		},
		Typed: typed,
	}, nil
}
