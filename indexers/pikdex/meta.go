package pikdex

import (
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
		s.Icon = string([]rune(strip(content))[0:2])
	},
}

func strip(input string) string {
	return strings.TrimSpace(input)
}
