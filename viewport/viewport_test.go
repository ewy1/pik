//go:build test

package viewport

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCrop(t *testing.T) {
	input := `0000
AAAA
` + Caret + `BBB
CCC`
	expected := `AAAA
` + Caret + `BBB`
	result, _, _ := Crop(input, strings.Split(input, "\n"), 2)
	assert.Equal(t, expected, result)
}

func TestCrop_Under(t *testing.T) {
	input := `0000
AAAA
` + Caret + `BBB`
	expected := `AAAA
` + Caret + `BBB`
	result, _, _ := Crop(input, strings.Split(input, "\n"), 2)
	assert.Equal(t, expected, result)
}

func TestCrop_Unnecessary(t *testing.T) {
	input := `AAAA
` + Caret + `BBB
CCC
DDDD`
	expected := input
	result, _, _ := Crop(input, strings.Split(input, "\n"), 8)
	assert.Equal(t, expected, result)
}

func TestNeedsViewport(t *testing.T) {
	amount := 3
	input := strings.Repeat("\n", amount)
	output := NeedsViewport(input, 4)
	assert.Equal(t, false, output)
}

func TestNeedsViewport_Negative(t *testing.T) {
	amount := 8
	input := strings.Repeat("\n", amount)
	output := NeedsViewport(input, 4)
	assert.Equal(t, true, output)
}
