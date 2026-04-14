package menu

import (
	"github.com/charmbracelet/lipgloss"
	"pik/menu/style"
	"pik/model"
)

var (
	SourceStyle = style.New(func() lipgloss.Style {
		st := lipgloss.NewStyle()
		return st
	})
	SourceHeaderBackground = lipgloss.Color("5")
	SourceHeaderStyle      = style.New(func() lipgloss.Style {
		st := lipgloss.NewStyle()
		return st

	})
	SourceLabelStyle = style.New(func() lipgloss.Style {
		st := lipgloss.NewStyle().Border(lipgloss.OuterHalfBlockBorder(), false, false, false, true).Background(SourceHeaderBackground).BorderBackground(SourceHeaderBackground).PaddingRight(1).PaddingLeft(1).MarginRight(1)
		return st

	})
	SourceIconStyle = style.New(func() lipgloss.Style {
		st := lipgloss.NewStyle().PaddingRight(1)
		return st

	})
	SourceTargetsStyle = style.New(func() lipgloss.Style {
		st := lipgloss.NewStyle()
		return st
	})
	SourcePathStyle = style.New(func() lipgloss.Style {
		st := lipgloss.NewStyle().Faint(true)
		return st
	})
)

func (m *Model) Source(src *model.HydratedSource) string {
	targets := make([]string, 0, len(src.Targets))
	for _, t := range src.HydratedTargets {
		targets = append(targets, m.Target(t))
	}

	targetContent := lipgloss.JoinVertical(lipgloss.Top, targets...)

	icon := SourceIconStyle.Render(Icon(src.Icon))

	return SourceStyle.Render(lipgloss.JoinVertical(lipgloss.Top,
		SourceHeaderStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, SourceLabelStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, icon, src.Label()), SourcePathStyle.Render(src.ShortPath())))),
		SourceTargetsStyle.Render(targetContent),
	))

}

var (
	StateStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle().MarginBottom(1)
	})
)

func (m *Model) State(st *model.HydratedState) string {
	sources := make([]string, 0, len(st.Sources))
	for _, hs := range st.HydratedSources {
		sources = append(sources, m.Source(hs))
	}

	return StateStyle.Render(lipgloss.JoinVertical(lipgloss.Top, sources...))
}
