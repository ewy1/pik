package completion

import (
	_ "embed"
	"github.com/ewy1/pik/spool"
)

//go:embed completion.sh
var completionCode string

var completionCodeByShell = map[string]string{
	"bash": ". <(pik --completion)",
	"zsh":  `autoload bashcompinit && bashcompinit && source <(pik --completion)`,
}

func Echo() error {
	_, err := spool.Print(completionCode)
	return err
}
