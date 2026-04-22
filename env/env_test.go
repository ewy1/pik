//go:build test

package env

import (
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"pik/model"
	"testing"
)

func TestIsEnv(t *testing.T) {
	err := pflag.Set("env", "test")
	assert.NoError(t, err)
	files := []string{
		".env",
		".env-test",
		".env.test",
		"test.env",
	}
	fakes := []string{
		".env-other",
		"unrelated.env",
	}
	for _, f := range files {
		assert.True(t, IsEnv(f))
	}
	for _, f := range fakes {
		assert.False(t, IsEnv(f))
	}
}

func TestGet(t *testing.T) {
	err := pflag.Set("env", "get")
	dir := t.TempDir()
	ctx := filepath.Join(dir, ".pik")
	err = os.Mkdir(ctx, 0777)
	assert.NoError(t, err)
	file := filepath.Join(dir, ".get.env")
	ctxFile := filepath.Join(dir, ".pik", ".get.env")
	src := &model.Source{
		Path: dir,
	}
	err = os.WriteFile(file, []byte("OUTERKEY=val"), 0666)
	assert.NoError(t, err)
	err = os.WriteFile(ctxFile, []byte("INNERKEY=otherval"), 0666)
	assert.NoError(t, err)
	res := Get(src)
	assert.Len(t, res, 2)
	assert.Contains(t, res, "OUTERKEY=val")
	assert.Contains(t, res, "INNERKEY=otherval")
}
