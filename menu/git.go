package menu

import (
	"github.com/charmbracelet/lipgloss"
	"pik/menu/style"
	"pik/model"
	"strconv"
)

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
