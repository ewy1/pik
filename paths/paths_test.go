//go:build test

package paths

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestSetAndReset(t *testing.T) {
	for i, pt := range Paths {
		old := pt.String()
		val := strconv.Itoa(i)
		Set(pt, val)
		assert.Equal(t, pt.String(), val)
		Reset()
		assert.Equal(t, pt.String(), old)
	}
}
