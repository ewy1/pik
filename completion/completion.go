package completion

import (
	_ "embed"
	"github.com/ewy1/pik/spool"
)

//go:embed completion.sh
var completionCode string

var completionCodeByShell = map[string]string{
	"bash": ". <(pik --completion)",
	"zsh":  `autoload -Uz compinit && compinit && source <(pik --completion)`,
}

func Echo() error {
	_, err := spool.Print(completionCode)
	return err
}
