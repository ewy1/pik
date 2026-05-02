package menu

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) HandleInput(msg tea.KeyMsg) (tea.Cmd, error) {

	if m.Search.Focused() {
		var cmd tea.Cmd
		switch msg.String() {
		case "ctrl+c":
			m.Search.SetValue("")
			m.Search.Blur()
		case "ctrl+d":
			m.Search.Blur()
		case "enter":
			m.Search.Blur()
		default:
			result, c := m.Search.Update(msg)
			cmd = c
			m.Search = result
		}
		return cmd, nil
	}

	var cmd tea.Cmd
	switch msg.String() {
	case "/":
		m.Search.SetValue("")
		fallthrough
	case "?":
		return m.Search.Focus(), nil
	case "i", "I":
		if m.Alt {
			m.Alt = false
			m.AutoAlt = false
			return tea.ExitAltScreen, nil
		} else {
			m.Alt = true
			m.AutoAlt = false
			return tea.EnterAltScreen, nil
		}
	case "h", "left":
		m.Leap(-1)
	case "l", "right":
		m.Leap(1)
	case "up", "k":
		m.Index--
	case "down", "j":
		m.Index++
	case "n":
		m.LeapFilter(1)
	case "N":
		m.LeapFilter(-1)
	case "q", "esc", "ctrl+c":
		m.Cancel = true
		return tea.Quit, nil
	case "space", " ", "enter", "ctrl+d":
		m.Done = true
		return tea.Quit, nil
	}

	_ = m.Validate()

	return cmd, nil
}

func (m *Model) LeapFilter(direction int) {
	startIndex := m.Index
	for {
		m.Index += direction
		clamped := m.Validate()
		if clamped {
			m.Index = startIndex
			return
		}

		source, target := m.Result()
		if m.Highlights(source, target) {
			return
		}
	}
}

func (m *Model) Leap(direction int) {
	for {
		source, target := m.Result()
		m.Index += direction
		m.Validate()
		newSource, newTarget := m.Result()
		if target == newTarget {
			return
		}

		if source != newSource {
			return
		}
	}
}
