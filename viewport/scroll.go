package viewport

import (
	"github.com/charmbracelet/lipgloss"
	"pik/menu/style"
	"strings"
)

var (
	ScrollTop          = StyleBarBackground.Render("╷")
	ScrollSpace        = StyleBarBackground.Render("│")
	ScrollBarTopEnd    = StyleBar.Render("╓")
	ScrollBar          = StyleBar.Render("║")
	ScrollBarBottomEnd = StyleBar.Render("╙")
	ScrollBottom       = StyleBarBackground.Render("╵")
)

var (
	StyleBar = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle()
	})
	StyleBarBackground = style.New(func() lipgloss.Style {
		return StyleBar.Get().Faint(true)
	})
)

func WithScroll(input string, barBegin int, barEnd int) string {
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		selection := ScrollSpace
		switch {
		case i == barBegin:
			selection = StyleBar.Render(ScrollBarTopEnd)
		case i == 0:
			selection = StyleBarBackground.Render(ScrollTop)
		case i == barEnd:
			selection = StyleBar.Render(ScrollBarBottomEnd)
		case i > barBegin && i < barEnd:
			selection = StyleBar.Render(ScrollBar)
		case i == len(lines)-1:
			selection = StyleBar.Render(ScrollBottom)
		}
		lines[i] = selection + " " + line
	}
	return strings.Join(lines, "\n")
}
