package menu

import (
	"github.com/charmbracelet/lipgloss"
	"pik/menu/style"
)

var (
	StateStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle().MarginBottom(1)
	})
	MotdStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle().Faint(true).PaddingLeft(1)
	})
)

func (m *Model) State() string {
	st := m.HydratedState
	sources := make([]string, 0, len(st.Sources))
	for i, hs := range st.HydratedSources {
		src := m.Source(hs)
		if i != len(st.HydratedSources)-1 {
			src += "\n"
		}
		sources = append(sources, src)
	}

	return StateStyle.Render(lipgloss.JoinVertical(lipgloss.Top, sources...), MotdStyle.Render("\n \U000F08B7  "+m.Motd))
}
