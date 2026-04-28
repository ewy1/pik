package menu

import (
	"github.com/charmbracelet/lipgloss"
	"pik/menu/style"
	"pik/model"
	"slices"
	"strings"
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
	var sub []string
	for _, t := range src.HydratedTargets {
		ts := t.Sub()
		header := !slices.Equal(sub, ts)
		if header {
			sub = ts
		}
		if header && strings.Join(ts, " ") != t.ShortestId() {
			targets = append(targets, m.Category(strings.Join(ts, " "), ""))
			header = false
		}
		targets = append(targets, m.Target(t, header))
	}

	targetContent := lipgloss.JoinVertical(lipgloss.Top, targets...)

	icon := PaddedIcon(src.Icon)

	parts := []string{
		SourceHeaderStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, SourceLabelStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, icon, src.Label()), SourcePathStyle.Render(src.ShortPath())))),
		SourceTargetsStyle.Render(targetContent),
	}

	if src.Git != nil {
		parts = append(parts, Git(src.Git))
	}

	return SourceStyle.Render(lipgloss.JoinVertical(lipgloss.Top,
		parts...,
	))
}
