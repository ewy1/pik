package style

import "github.com/charmbracelet/lipgloss"

func (s Style) Debug() Style {
	return New(func() lipgloss.Style {
		return s.builder().Foreground(lipgloss.Color("2")).Background(lipgloss.Color("1"))
	})
}
