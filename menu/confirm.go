package menu

import (
	"bufio"
	"fmt"
	"git.ewy.one/pik.git/flags"
	"git.ewy.one/pik.git/menu/style"
	"git.ewy.one/pik.git/model"
	"github.com/charmbracelet/lipgloss"
	"io"
	"os"
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
		st := lipgloss.NewStyle()
		if !*flags.Yes {
			st.Foreground(PromptColor)
		}
		return st
	})
	ConfirmStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle()
	})
	AnswerStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle()
	})
	YesStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle().Faint(true)
	})
)

func Confirm(r io.Reader, source *model.Source, target model.Target, args ...string) bool {
	parts := []string{
		ConfirmStyle.Render(PromptStyle.Render("[Y/n]")),
		Banner(source, target, args...),
		"? ",
	}
	banner := BannerStyle.Render(parts...)
	_, _ = fmt.Fprint(os.Stderr, banner)

	if *flags.Yes {
		_, _ = fmt.Fprintln(os.Stderr, AnswerStyle.Render("Y", YesStyle.Render("(--yes)")))
		return true
	}

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
