package shell

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShell_ShellByFilename(t *testing.T) {
	s := &shell{
		Locations: map[string]string{"powershell": "/pws", "bash": "/bash"},
	}
	inputs := map[string]string{
		"script.ps1": "/pws",
		"script.sh":  "/bash",
		"asdf":       "",
	}

	for k, v := range inputs {
		sh := s.ShellByFilename(k)
		assert.Equal(t, sh, v)
	}
}
