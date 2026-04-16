package menu

import (
	"github.com/charmbracelet/lipgloss"
	"pik/menu/style"
	"pik/model"
	"strconv"
)

var (
	SourceStyle = style.New(func() lipgloss.Style {
		st := lipgloss.NewStyle().PaddingBottom(1)
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
		st := lipgloss.NewStyle().PaddingRight(1).Width(3)
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

	icon := SourceIconStyle.Render(src.Icon)

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

var (
	GitColor     = lipgloss.Color("4")
	GitInfoStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle().Background(GitColor).Border(lipgloss.OuterHalfBlockBorder(), false, false, false, true).BorderBackground(GitColor).Padding(0, 1)
	})
	GitStatusStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle().Bold(true).Background(GitColor).PaddingLeft(1)
	})
	GitAddColor = lipgloss.Color("2")
	GitAddStyle = style.New(func() lipgloss.Style {
		return GitStatusStyle.Get().Foreground(GitAddColor)
	})
	GitRemoveColor = lipgloss.Color("1")
	GitRemoveStyle = style.New(func() lipgloss.Style {
		return GitStatusStyle.Get().Foreground(GitRemoveColor)
	})
	GitChangeColor = lipgloss.Color("5")
	GitChangeStyle = style.New(func() lipgloss.Style {
		return GitStatusStyle.Get().Foreground(GitChangeColor)
	})
)

func Git(info *model.GitInfo) string {
	var parts = []string{
		"  ",
		info.Branch,
	}

	if info.Insertions > 0 {
		parts = append(parts, GitAddStyle.Render("+"+strconv.Itoa(info.Insertions)))
	}
	if info.Deletions > 0 {
		parts = append(parts, GitRemoveStyle.Render("-"+strconv.Itoa(info.Deletions)))
	}
	if info.Changes > 0 {
		parts = append(parts, GitChangeStyle.Render("~"+strconv.Itoa(info.Changes)))
	}
	if info.Changes == 0 && info.Deletions == 0 && info.Insertions == 0 {
		parts = append(parts, GitAddStyle.Render("clean"))
	}

	return GitInfoStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left,
		parts...,
	))
}
