package menu

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ewy1/pik/menu/style"
	"strings"
)

var (
	IconStyle = style.New(func() lipgloss.Style {
		st := lipgloss.NewStyle().Width(2).Height(1)
		return st
	})
)

func Icon(input string) string {
	return IconStyle.Render(input)
}

func PaddedIcon(input string) string {
	if strings.TrimSpace(input) == "" {
		return Icon("  ")
	}
	return Icon(input)
}
