package menu

import "strings"

func (m *Model) AddSearch(croppedInput string) string {
	lines := strings.Split(croppedInput, "\n")
	lastIndex := len(lines) - 1
	view := m.Search.View()
	lines[lastIndex] = view
	return strings.Join(lines, "\n")
}
