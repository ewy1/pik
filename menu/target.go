package menu

import (
	"github.com/charmbracelet/lipgloss"
	"pik/menu/style"
	"pik/model"
	"pik/viewport"
)

var (
	TargetBackgroundColor         = lipgloss.Color("8")
	SelectedTargetBackgroundColor = lipgloss.Color("2")
	TargetStyle                   = style.New(func() lipgloss.Style {
		st := lipgloss.NewStyle().Border(lipgloss.OuterHalfBlockBorder(), false, false, false, true)
		return st
	})
	SelectedTargetStyle = style.New(func() lipgloss.Style {
		return TargetStyle.Get().BorderBackground(SelectedTargetBackgroundColor).Background(SelectedTargetBackgroundColor)
	})
	TargetLabelStyle = style.New(func() lipgloss.Style {
		st := lipgloss.NewStyle().MarginRight(1)
		return st
	})
	TargetDescriptionStyle = style.New(func() lipgloss.Style {
		st := lipgloss.NewStyle().Faint(true).MarginLeft(1)
		return st
	})
	SelectedTargetDescriptionStyle = style.New(func() lipgloss.Style {
		st := TargetDescriptionStyle.Get().Faint(false)
		return st
	})
	TargetIconStyle = style.New(func() lipgloss.Style {
		st := lipgloss.NewStyle().PaddingLeft(1)
		return st
	})
	TargetIconSelectedStyle = style.New(func() lipgloss.Style {
		return TargetIconStyle.Get().MarginLeft(1).PaddingLeft(0)
	})
	TargetSubStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle()
	})
	CategoryColor = lipgloss.Color("7")
	CategoryStyle = style.New(func() lipgloss.Style {
		return TargetStyle.Get().BorderForeground(CategoryColor).Background(CategoryColor).BorderBackground(CategoryColor)
	})
	OrphanCategoryStyle = style.New(func() lipgloss.Style {
		return TargetLabelStyle.Get().PaddingLeft(1)
	})
)

func (m *Model) Target(t model.HydratedTarget, header bool) string {
	_, selection := m.Result()
	selected := selection != nil && selection.Target() == t.Target()
	icon := ""
	if selected {
		icon = TargetIconSelectedStyle.Render(PaddedIcon(viewport.Caret))
	} else {
		icon = TargetIconStyle.Render(PaddedIcon(t.Icon()))
	}
	selectionStyle := TargetStyle
	selectionDescriptionStyle := TargetDescriptionStyle
	if selected {
		selectionStyle = SelectedTargetStyle
		selectionDescriptionStyle = SelectedTargetDescriptionStyle
	} else if header {
		selectionStyle = CategoryStyle
	}
	var labelParts []string
	labelParts = append(labelParts, icon)
	sub := t.Sub()
	if sub != nil && sub[len(sub)-1] != t.ShortestId() {
		labelParts = append(labelParts, TargetSubStyle.Render(sub...))
	}
	labelParts = append(labelParts, TargetLabelStyle.Render(t.Label()))
	return lipgloss.JoinHorizontal(lipgloss.Left, selectionStyle.Render(labelParts...), selectionDescriptionStyle.Render(t.Description()))
}

func (m *Model) Category(input string, desc string) string {
	return CategoryStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, OrphanCategoryStyle.Render(input), TargetDescriptionStyle.Render(desc)))
}
