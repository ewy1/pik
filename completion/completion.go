package completion

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/ewy1/pik/paths"
	"github.com/ewy1/pik/spool"
	"os"
	"path/filepath"
	"strings"
)

//go:embed completion.sh
var completionCode string

var completionFormat = `
# %s
%s
`

var completionComment = "\n# pik completion (installed by `pik --install-completion`)\n"

var completionCodeByShell = map[string]string{
	"bash": ". <(pik --completion)\n\n",
	"zsh":  `autoload -Uz compinit && compinit && source <(pik --completion)\n\n`,
}

var completionFileByShell = map[string]string{
	"bash": ".bashrc",
	"zsh":  ".zshrc",
}

var AlreadyInstalledError = errors.New("completion already installed")

func Add(shell string) error {
	f := filepath.Join(paths.HomeDir.String(), completionFileByShell[shell])
	content, err := os.ReadFile(f)
	if err != nil {
		return err
	}
	if strings.Contains(string(content), strings.TrimSpace(completionCodeByShell[shell])) {
		return AlreadyInstalledError
	}
	fd, err := os.OpenFile(f, os.O_APPEND|os.O_WRONLY, 0600)
	defer fd.Close()
	if err != nil {
		return err
	}
	_, err = fd.Write([]byte(fmt.Sprintf(completionFormat, completionComment, completionCodeByShell[shell])))
	if err != nil {
		return err
	}
	successMessage(shell, f)
	return nil
}

func successMessage(shell string, file string) {
	_, _ = spool.Print("Installed completion for %s in %s\n", shell, file)
}

func Echo() error {
	_, err := spool.Print(completionCode)
	return err
}
