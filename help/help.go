package help

import (
	_ "embed"
	"github.com/ewy1/pik/spool"
)

//go:embed help.txt
var content string

func Echo() {
	_, _ = spool.Print("%v", content)
}
