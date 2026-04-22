//go:build test

package describe

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"pik/runner"
	"testing"
)

type key struct {
	*runner.Stub
}

func Key() *key {
	return &key{
		Stub: &runner.Stub{},
	}
}

func TestDescribe(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file")
	err := os.WriteFile(f, []byte("# this is a comment"), 0666)
	assert.NoError(t, err)
	result, err := Describe(Key(), f)
	assert.NoError(t, err)
	assert.Contains(t, result, "this is a comment")
}

func TestDescribe_Slashes(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file")
	err := os.WriteFile(f, []byte(" // this is a comment"), 0666)
	assert.NoError(t, err)
	result, err := Describe(Key(), f)
	assert.NoError(t, err)
	assert.Contains(t, result, "this is a comment")
}

func TestDescribe_Cache(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file")
	k := Key()
	err := os.WriteFile(f, []byte(" // this is a comment"), 0666)
	assert.NoError(t, err)
	result, err := Describe(k, f)
	assert.NoError(t, err)
	assert.Contains(t, result, "this is a comment")
	err = os.Remove(f)
	if err != nil {
		assert.NoError(t, err)
	}
	cached, err := Describe(k, f)
	assert.NoError(t, err)
	assert.Equal(t, result, cached)

}

func TestDescribe_NoDescription(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file")
	err := os.WriteFile(f, []byte("1230987324091823740918237409283"), 0666)
	assert.NoError(t, err)
	result, err := Describe(Key(), f)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestDescribe_Shebang_NoDescription(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file")
	err := os.WriteFile(f, []byte("#!/usr/bin/env bash\nkjausdhuahsdf"), 0666)
	assert.NoError(t, err)
	result, err := Describe(Key(), f)
	assert.Empty(t, result)
	assert.NoError(t, err)
}

func TestDescribe_Error(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file")
	result, err := Describe(Key(), f)
	assert.Empty(t, result)
	assert.Error(t, err)
}
