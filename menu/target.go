package menu

import (
	"github.com/charmbracelet/lipgloss"
	"pik/menu/style"
	"pik/model"
)

var (
	TargetBackgroundColor         = lipgloss.Color("8")
	SelectedTargetBackgroundColor = lipgloss.Color("2")
	TargetStyle                   = style.New(func() lipgloss.Style {
		st := lipgloss.NewStyle().Border(lipgloss.OuterHalfBlockBorder(), false, false, false, true).BorderBackground(TargetBackgroundColor).Background(TargetBackgroundColor)
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
	TargetSubStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle()
	})
)

func (m *Model) Target(t model.HydratedTarget) string {
	icon := TargetIconStyle.Render(PaddedIcon(t.Icon()))
	selectionStyle := TargetStyle
	selectionDescriptionStyle := TargetDescriptionStyle
	_, sel := m.Result()
	if sel.Target() == t.Target() {
		selectionStyle = SelectedTargetStyle
		selectionDescriptionStyle = SelectedTargetDescriptionStyle
	}
	var labelParts []string
	labelParts = append(labelParts, icon)
	if t.Sub() != nil {
		labelParts = append(labelParts, TargetSubStyle.Render(t.Sub()...))
	}
	labelParts = append(labelParts, TargetLabelStyle.Render(t.Label()))
	return lipgloss.JoinHorizontal(lipgloss.Left, selectionStyle.Render(labelParts...), selectionDescriptionStyle.Render(t.Description()))
}
