//go:build test

package paths

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var bHome, bThis, bCache, bConfig string

// Set temporarily sets the paths for unit test purposes
// remember to defer Reset
func Set(home, this, cache, config string) {
	bHome = Home
	bThis = This
	bCache = Cache
	bConfig = Config

	Home = home
	This = this
	Cache = cache
	Config = config
}

// Reset sets the path variables back to before the unit test
func Reset() {
	Home = bHome
	This = bThis
	Cache = bCache
	Config = bConfig
}

func TestSetAndReset(t *testing.T) {
	Set("1", "2", "3", "4")
	assert.Equal(t, Home, "1")
	assert.Equal(t, This, "2")
	assert.Equal(t, Cache, "3")
	assert.Equal(t, Config, "4")
	Reset()
	assert.NotEqual(t, Home, "1")
	assert.NotEqual(t, This, "2")
	assert.NotEqual(t, Cache, "3")
	assert.NotEqual(t, Config, "4")
}
