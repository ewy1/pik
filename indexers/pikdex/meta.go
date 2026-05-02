package pikdex

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type MetaSetter func(s *SourceData, content string)

var MetaFiles = map[string]MetaSetter{
	".alias": func(s *SourceData, content string) {
		split := strings.Split(content, "\n")
		s.Aliases = make([]string, 0, len(split))
		for _, line := range split {
			stripped := strip(line)
			if stripped != "" {
				s.Aliases = append(s.Aliases, stripped)
			}
		}
	},
	".icon": func(s *SourceData, content string) {
		icon := strip(content)
		desiredWidth := lipgloss.Width(icon)
		diff := desiredWidth - len([]rune(icon))
		icon += strings.Repeat(" ", diff)
		s.Icon = icon
	},
}

func strip(input string) string {
	return strings.TrimSpace(input)
}
