package menu

import (
	"bufio"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"io"
	"os"
	"pik/menu/style"
	"pik/model"
	"slices"
)

var confirmations = []rune{
	'y',
	'Y',
	' ',
	'\n',
}

var (
	PromptColor = lipgloss.Color("1")
	PromptStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle().Foreground(PromptColor)
	})
	ConfirmStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle()
	})
)

func Confirm(r io.Reader, source *model.Source, target model.Target, args ...string) bool {
	banner := ConfirmStyle.Render(PromptStyle.Render("[Y/n]"), Banner(source, target, args...), "? ")
	fmt.Print(banner)
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanRunes)
	scanner.Scan()
	if slices.Contains(confirmations, []rune(scanner.Text())[0]) {
		return true
	} else {
		_, _ = fmt.Fprint(os.Stderr, "confirmation was not given.")
	}
	return false
}
