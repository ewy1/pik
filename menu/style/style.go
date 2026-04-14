package style

import "github.com/charmbracelet/lipgloss"

type StyleBuilder func() lipgloss.Style

type Style struct {
	style   *lipgloss.Style
	builder StyleBuilder
}

func New(builder StyleBuilder) Style {
	return Style{
		builder: builder,
	}
}

func (s *Style) Get() lipgloss.Style {

	if s.style == nil {
		st := s.builder()
		s.style = &st
	}

	return *s.style
}

func (s *Style) Render(input ...string) string {
	return s.Get().Render(input...)
}
